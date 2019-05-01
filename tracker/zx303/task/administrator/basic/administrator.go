package basic

import (
	"fmt"
	brainException "gitlab.com/iotTracker/brain/exception"
	zx303Task "gitlab.com/iotTracker/brain/tracker/zx303/task"
	zx303TaskAction "gitlab.com/iotTracker/brain/tracker/zx303/task/action"
	zx303TaskAdministrator "gitlab.com/iotTracker/brain/tracker/zx303/task/administrator"
	zx303TaskAdministratorException "gitlab.com/iotTracker/brain/tracker/zx303/task/administrator/exception"
	zx303TaskRecordHandler "gitlab.com/iotTracker/brain/tracker/zx303/task/recordHandler"
	zx303TaskStep "gitlab.com/iotTracker/brain/tracker/zx303/task/step"
	zx303TaskValidator "gitlab.com/iotTracker/brain/tracker/zx303/task/validator"
	zx303TaskSubmittedMessage "gitlab.com/iotTracker/messaging/message/zx303/task/submitted"
	messagingProducer "gitlab.com/iotTracker/messaging/producer"
)

type administrator struct {
	zx303TaskValidator     zx303TaskValidator.Validator
	zx303TaskRecordHandler *zx303TaskRecordHandler.RecordHandler
	nerveBroadcastProducer messagingProducer.Producer
}

func New(
	zx303TaskValidator zx303TaskValidator.Validator,
	zx303TaskRecordHandler *zx303TaskRecordHandler.RecordHandler,
	nerveBroadcastProducer messagingProducer.Producer,
) zx303TaskAdministrator.Administrator {
	return &administrator{
		zx303TaskValidator:     zx303TaskValidator,
		zx303TaskRecordHandler: zx303TaskRecordHandler,
		nerveBroadcastProducer: nerveBroadcastProducer,
	}
}

func (a *administrator) ValidateSubmitRequest(request *zx303TaskAdministrator.SubmitRequest) error {
	reasonsInvalid := make([]string, 0)

	if request.Claims == nil {
		reasonsInvalid = append(reasonsInvalid, "claims are nil")
	} else {
		zx303DeviceValidateResponse, err := a.zx303TaskValidator.Validate(&zx303TaskValidator.ValidateRequest{
			Claims:    request.Claims,
			ZX303Task: request.ZX303Task,
			Action:    zx303TaskAction.Create,
		})
		if err != nil {
			reasonsInvalid = append(reasonsInvalid, "error validating zx303 task: "+err.Error())
		}
		if len(zx303DeviceValidateResponse.ReasonsInvalid) > 0 {
			for _, reason := range zx303DeviceValidateResponse.ReasonsInvalid {
				reasonsInvalid = append(reasonsInvalid, fmt.Sprintf("zx303 task invalid: %s - %s - %s", reason.Field, reason.Type, reason.Help))
			}
		}
	}

	if len(reasonsInvalid) > 0 {
		return brainException.RequestInvalid{Reasons: reasonsInvalid}
	}

	return nil
}

func (a *administrator) Submit(request *zx303TaskAdministrator.SubmitRequest) (*zx303TaskAdministrator.SubmitResponse, error) {
	if err := a.ValidateSubmitRequest(request); err != nil {
		return nil, err
	}

	// create task
	createResponse, err := a.zx303TaskRecordHandler.Create(&zx303TaskRecordHandler.CreateRequest{
		ZX303Task: request.ZX303Task,
	})
	if err != nil {
		return nil, zx303TaskAdministratorException.ZX303TaskSubmission{Reasons: []string{"creation", err.Error()}}
	}

	// produce task generated event to nerveBroadcast topic
	if err := a.nerveBroadcastProducer.Produce(zx303TaskSubmittedMessage.Message{
		Task: createResponse.ZX303Task,
	}); err != nil {
		return nil, zx303TaskAdministratorException.ZX303TaskSubmission{Reasons: []string{"message production", err.Error()}}
	}

	return &zx303TaskAdministrator.SubmitResponse{
		ZX303Task: createResponse.ZX303Task,
	}, nil
}

func (a *administrator) ValidateFailTaskRequest(request *zx303TaskAdministrator.FailTaskRequest) error {
	reasonsInvalid := make([]string, 0)

	if request.Claims == nil {
		reasonsInvalid = append(reasonsInvalid, "claims are nil")
	}

	if request.ZX303TaskIdentifier == nil {
		reasonsInvalid = append(reasonsInvalid, "zx303TaskIdentifier is nil")
	}

	if len(reasonsInvalid) > 0 {
		return brainException.RequestInvalid{Reasons: reasonsInvalid}
	}

	return nil
}

func (a *administrator) FailTask(request *zx303TaskAdministrator.FailTaskRequest) (*zx303TaskAdministrator.FailTaskResponse, error) {
	if err := a.ValidateFailTaskRequest(request); err != nil {
		return nil, err
	}

	// retrieve the task
	retrieveResponse, err := a.zx303TaskRecordHandler.Retrieve(&zx303TaskRecordHandler.RetrieveRequest{
		Claims:     request.Claims,
		Identifier: request.ZX303TaskIdentifier,
	})
	if err != nil {
		return nil, zx303TaskAdministratorException.ZX303TaskFail{Reasons: []string{"retrieval", err.Error()}}
	}

	// change task status to failed
	retrieveResponse.ZX303Task.Status = zx303Task.Failed

	// change step status to failed if step idx provided
	if request.FailedStepIdx >= 0 {
		if request.FailedStepIdx > len(retrieveResponse.ZX303Task.Steps)-1 {
			return nil, zx303TaskAdministratorException.ZX303TaskFail{Reasons: []string{"invalid step idx", string(request.FailedStepIdx)}}

		}
		retrieveResponse.ZX303Task.Steps[request.FailedStepIdx].Status = zx303TaskStep.Failed
	}

	// update task
	if _, err := a.zx303TaskRecordHandler.Update(&zx303TaskRecordHandler.UpdateRequest{
		Claims:     request.Claims,
		Identifier: request.ZX303TaskIdentifier,
		ZX303Task:  retrieveResponse.ZX303Task,
	}); err != nil {
		return nil, zx303TaskAdministratorException.ZX303TaskFail{Reasons: []string{"update", err.Error()}}
	}

	return &zx303TaskAdministrator.FailTaskResponse{
		ZX303Task: retrieveResponse.ZX303Task,
	}, nil
}

func (a *administrator) ValidateTransitionTaskRequest(request *zx303TaskAdministrator.TransitionTaskRequest) error {
	reasonsInvalid := make([]string, 0)

	if request.Claims == nil {
		reasonsInvalid = append(reasonsInvalid, "claims are nil")
	}

	if request.ZX303TaskIdentifier == nil {
		reasonsInvalid = append(reasonsInvalid, "zx303TaskIdentifier is nil")
	}

	if len(reasonsInvalid) > 0 {
		return brainException.RequestInvalid{Reasons: reasonsInvalid}
	}
	return nil
}

func (a *administrator) TransitionTask(request *zx303TaskAdministrator.TransitionTaskRequest) (*zx303TaskAdministrator.TransitionTaskResponse, error) {
	if err := a.ValidateTransitionTaskRequest(request); err != nil {
		return nil, err
	}

	// retrieve the task
	retrieveResponse, err := a.zx303TaskRecordHandler.Retrieve(&zx303TaskRecordHandler.RetrieveRequest{
		Claims:     request.Claims,
		Identifier: request.ZX303TaskIdentifier,
	})
	if err != nil {
		return nil, zx303TaskAdministratorException.ZX303TaskTransition{Reasons: []string{"retrieval", err.Error()}}
	}

	if request.StepIdx > len(retrieveResponse.ZX303Task.Steps)-1 {
		return nil, zx303TaskAdministratorException.ZX303TaskTransition{Reasons: []string{"invalid step idx", string(request.StepIdx)}}
	}

	// update the step status
	retrieveResponse.ZX303Task.Steps[request.StepIdx].Status = request.NewStepStatus

	// check how the task should transition
	if request.StepIdx < len(retrieveResponse.ZX303Task.Steps)-1 {
		// this is not the last step
		// if the task status is not yet in executing update it
		if retrieveResponse.ZX303Task.Status != zx303Task.Executing {
			retrieveResponse.ZX303Task.Status = zx303Task.Executing
		}
	} else {
		// this is the last step
		// if the new step status is finished the task is finished
		if request.NewStepStatus == zx303TaskStep.Finished {
			retrieveResponse.ZX303Task.Status = zx303Task.Finished
		}
	}

	// update the task
	if _, err := a.zx303TaskRecordHandler.Update(&zx303TaskRecordHandler.UpdateRequest{
		Claims:     request.Claims,
		Identifier: request.ZX303TaskIdentifier,
		ZX303Task:  retrieveResponse.ZX303Task,
	}); err != nil {
		return nil, zx303TaskAdministratorException.ZX303TaskFail{Reasons: []string{"update", err.Error()}}
	}

	return &zx303TaskAdministrator.TransitionTaskResponse{
		ZX303Task: retrieveResponse.ZX303Task,
	}, nil
}

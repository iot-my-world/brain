package basic

import (
	"fmt"
	brainException "github.com/iot-my-world/brain/internal/exception"
	"github.com/iot-my-world/brain/pkg/tracker/zx303/task"
	"github.com/iot-my-world/brain/pkg/tracker/zx303/task/action"
	administrator2 "github.com/iot-my-world/brain/pkg/tracker/zx303/task/administrator"
	"github.com/iot-my-world/brain/pkg/tracker/zx303/task/administrator/exception"
	"github.com/iot-my-world/brain/pkg/tracker/zx303/task/recordHandler"
	"github.com/iot-my-world/brain/pkg/tracker/zx303/task/step"
	"github.com/iot-my-world/brain/pkg/tracker/zx303/task/validator"
)

type administrator struct {
	zx303TaskValidator     validator.Validator
	zx303TaskRecordHandler *recordHandler.RecordHandler
	//nerveBroadcastProducer messagingProducer.Producer
}

func New(
	zx303TaskValidator validator.Validator,
	zx303TaskRecordHandler *recordHandler.RecordHandler,
	//nerveBroadcastProducer messagingProducer.Producer,
) administrator2.Administrator {
	return &administrator{
		zx303TaskValidator:     zx303TaskValidator,
		zx303TaskRecordHandler: zx303TaskRecordHandler,
		//nerveBroadcastProducer: nerveBroadcastProducer,
	}
}

func (a *administrator) ValidateSubmitRequest(request *administrator2.SubmitRequest) error {
	reasonsInvalid := make([]string, 0)

	if request.Claims == nil {
		reasonsInvalid = append(reasonsInvalid, "claims are nil")
	} else {
		zx303DeviceValidateResponse, err := a.zx303TaskValidator.Validate(&validator.ValidateRequest{
			Claims:    request.Claims,
			ZX303Task: request.ZX303Task,
			Action:    action.Create,
		})
		if err != nil {
			reasonsInvalid = append(reasonsInvalid, "error validating zx303 task: "+err.Error())
		} else {
			if len(zx303DeviceValidateResponse.ReasonsInvalid) > 0 {
				for _, reason := range zx303DeviceValidateResponse.ReasonsInvalid {
					reasonsInvalid = append(reasonsInvalid, fmt.Sprintf("zx303 task invalid: %s - %s - %s", reason.Field, reason.Type, reason.Help))
				}
			}
		}
	}

	if len(reasonsInvalid) > 0 {
		return brainException.RequestInvalid{Reasons: reasonsInvalid}
	}

	return nil
}

func (a *administrator) Submit(request *administrator2.SubmitRequest) (*administrator2.SubmitResponse, error) {
	if err := a.ValidateSubmitRequest(request); err != nil {
		return nil, err
	}

	// create task
	createResponse, err := a.zx303TaskRecordHandler.Create(&recordHandler.CreateRequest{
		ZX303Task: request.ZX303Task,
	})
	if err != nil {
		return nil, exception.ZX303TaskSubmission{Reasons: []string{"creation", err.Error()}}
	}

	// produce task generated message to nerveBroadcast topic
	//if err := a.nerveBroadcastProducer.Produce(zx303TaskSubmittedMessage.Message{
	//	Task: createResponse.ZX303Task,
	//}); err != nil {
	//	return nil, zx303TaskAdministratorException.ZX303TaskSubmission{Reasons: []string{"message production", err.Error()}}
	//}

	return &administrator2.SubmitResponse{
		ZX303Task: createResponse.ZX303Task,
	}, nil
}

func (a *administrator) ValidateFailTaskRequest(request *administrator2.FailTaskRequest) error {
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

func (a *administrator) FailTask(request *administrator2.FailTaskRequest) (*administrator2.FailTaskResponse, error) {
	if err := a.ValidateFailTaskRequest(request); err != nil {
		return nil, err
	}

	// retrieve the task
	retrieveResponse, err := a.zx303TaskRecordHandler.Retrieve(&recordHandler.RetrieveRequest{
		Claims:     request.Claims,
		Identifier: request.ZX303TaskIdentifier,
	})
	if err != nil {
		return nil, exception.ZX303TaskFail{Reasons: []string{"retrieval", err.Error()}}
	}

	// change task status to failed
	retrieveResponse.ZX303Task.Status = task.Failed

	// change step status to failed if step idx provided
	if request.FailedStepIdx >= 0 {
		if request.FailedStepIdx > len(retrieveResponse.ZX303Task.Steps)-1 {
			return nil, exception.ZX303TaskFail{Reasons: []string{"invalid step idx", string(request.FailedStepIdx)}}

		}
		retrieveResponse.ZX303Task.Steps[request.FailedStepIdx].Status = step.Failed
	}

	// update task
	if _, err := a.zx303TaskRecordHandler.Update(&recordHandler.UpdateRequest{
		Claims:     request.Claims,
		Identifier: request.ZX303TaskIdentifier,
		ZX303Task:  retrieveResponse.ZX303Task,
	}); err != nil {
		return nil, exception.ZX303TaskFail{Reasons: []string{"update", err.Error()}}
	}

	return &administrator2.FailTaskResponse{
		ZX303Task: retrieveResponse.ZX303Task,
	}, nil
}

func (a *administrator) ValidateTransitionTaskRequest(request *administrator2.TransitionTaskRequest) error {
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

func (a *administrator) TransitionTask(request *administrator2.TransitionTaskRequest) (*administrator2.TransitionTaskResponse, error) {
	if err := a.ValidateTransitionTaskRequest(request); err != nil {
		return nil, err
	}

	// retrieve the task
	retrieveResponse, err := a.zx303TaskRecordHandler.Retrieve(&recordHandler.RetrieveRequest{
		Claims:     request.Claims,
		Identifier: request.ZX303TaskIdentifier,
	})
	if err != nil {
		return nil, exception.ZX303TaskTransition{Reasons: []string{"retrieval", err.Error()}}
	}

	if request.StepIdx > len(retrieveResponse.ZX303Task.Steps)-1 {
		return nil, exception.ZX303TaskTransition{Reasons: []string{"invalid step idx", string(request.StepIdx)}}
	}

	// update the step status
	retrieveResponse.ZX303Task.Steps[request.StepIdx].Status = request.NewStepStatus

	// check how the task should transition
	if request.StepIdx < len(retrieveResponse.ZX303Task.Steps)-1 {
		// this is not the last step
		// if the task status is not yet in executing update it
		if retrieveResponse.ZX303Task.Status != task.Executing {
			retrieveResponse.ZX303Task.Status = task.Executing
		}

		//// produce task transitioned message to nerveBroadcast topic
		//if err := a.nerveBroadcastProducer.Produce(zx303TaskTransitionedMessage.Message{
		//	Task: retrieveResponse.ZX303Task,
		//}); err != nil {
		//	return nil, zx303TaskAdministratorException.ZX303TaskTransition{Reasons: []string{"message production", err.Error()}}
		//}
	} else {
		// this is the last step
		// if the new step status is finished the task is finished
		if request.NewStepStatus == step.Finished {
			retrieveResponse.ZX303Task.Status = task.Finished
		}
	}

	// update the task
	if _, err := a.zx303TaskRecordHandler.Update(&recordHandler.UpdateRequest{
		Claims:     request.Claims,
		Identifier: request.ZX303TaskIdentifier,
		ZX303Task:  retrieveResponse.ZX303Task,
	}); err != nil {
		return nil, exception.ZX303TaskFail{Reasons: []string{"update", err.Error()}}
	}

	return &administrator2.TransitionTaskResponse{
		ZX303Task: retrieveResponse.ZX303Task,
	}, nil
}

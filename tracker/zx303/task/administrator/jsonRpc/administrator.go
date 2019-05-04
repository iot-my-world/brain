package jsonRpc

import (
	jsonRpcClient "gitlab.com/iotTracker/brain/communication/jsonRpc/client"
	brainException "gitlab.com/iotTracker/brain/exception"
	wrappedIdentifier "gitlab.com/iotTracker/brain/search/identifier/wrapped"
	zx303TaskAdministrator "gitlab.com/iotTracker/brain/tracker/zx303/task/administrator"
	zx303TaskAdministratorJsonRpcAdaptor "gitlab.com/iotTracker/brain/tracker/zx303/task/administrator/adaptor/jsonRpc"
)

type administrator struct {
	jsonRpcClient jsonRpcClient.Client
}

func New(
	jsonRpcClient jsonRpcClient.Client,
) zx303TaskAdministrator.Administrator {
	return &administrator{
		jsonRpcClient: jsonRpcClient,
	}
}

func (a *administrator) Submit(request *zx303TaskAdministrator.SubmitRequest) (*zx303TaskAdministrator.SubmitResponse, error) {
	return nil, brainException.NotImplemented{}
}

func (a *administrator) ValidateFailTaskRequest(request *zx303TaskAdministrator.FailTaskRequest) error {
	reasonsInvalid := make([]string, 0)

	if request.ZX303TaskIdentifier == nil {
		reasonsInvalid = append(reasonsInvalid, "identifier is nil")
	}

	if !a.jsonRpcClient.LoggedIn() {
		reasonsInvalid = append(reasonsInvalid, "json rpc client is not logged in")
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

	// create wrapped identifier
	wrappedZX303TaskIdentifier, err := wrappedIdentifier.Wrap(request.ZX303TaskIdentifier)
	if err != nil {
		return nil, brainException.Unexpected{Reasons: []string{"wrapping device identifier", err.Error()}}
	}

	// perform request
	failTaskResponse := zx303TaskAdministratorJsonRpcAdaptor.FailTaskResponse{}
	if err := a.jsonRpcClient.JsonRpcRequest(
		"ZX303TaskAdministrator.FailTask",
		zx303TaskAdministratorJsonRpcAdaptor.FailTaskRequest{
			ZX303TaskIdentifier: *wrappedZX303TaskIdentifier,
		},
		&failTaskResponse,
	); err != nil {
		return nil, brainException.Unexpected{Reasons: []string{"fail task error", err.Error()}}
	}

	return &zx303TaskAdministrator.FailTaskResponse{
		ZX303Task: failTaskResponse.ZX303Task,
	}, nil
}

func (a *administrator) ValidateTransitionTaskRequest(request *zx303TaskAdministrator.TransitionTaskRequest) error {
	reasonsInvalid := make([]string, 0)

	if request.ZX303TaskIdentifier == nil {
		reasonsInvalid = append(reasonsInvalid, "identifier is nil")
	}

	if !a.jsonRpcClient.LoggedIn() {
		reasonsInvalid = append(reasonsInvalid, "json rpc client is not logged in")
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

	// create wrapped identifier
	wrappedZX303TaskIdentifier, err := wrappedIdentifier.Wrap(request.ZX303TaskIdentifier)
	if err != nil {
		return nil, brainException.Unexpected{Reasons: []string{"wrapping device identifier", err.Error()}}
	}

	// perform request
	transitionTaskResponse := zx303TaskAdministratorJsonRpcAdaptor.TransitionTaskResponse{}
	if err := a.jsonRpcClient.JsonRpcRequest(
		"ZX303TaskAdministrator.TransitionTask",
		zx303TaskAdministratorJsonRpcAdaptor.TransitionTaskRequest{
			ZX303TaskIdentifier: *wrappedZX303TaskIdentifier,
			StepIdx:             request.StepIdx,
			NewStepStatus:       request.NewStepStatus,
		},
		&transitionTaskResponse,
	); err != nil {
		return nil, brainException.Unexpected{Reasons: []string{"fail task error", err.Error()}}
	}

	return &zx303TaskAdministrator.TransitionTaskResponse{
		ZX303Task: transitionTaskResponse.ZX303Task,
	}, nil
}

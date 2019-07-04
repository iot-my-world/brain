package jsonRpc

import (
	jsonRpcClient "github.com/iot-my-world/brain/communication/jsonRpc/client"
	brainException "github.com/iot-my-world/brain/exception"
	wrappedIdentifier "github.com/iot-my-world/brain/pkg/search/identifier/wrapped"
	administrator2 "github.com/iot-my-world/brain/pkg/tracker/zx303/task/administrator"
	"github.com/iot-my-world/brain/pkg/tracker/zx303/task/administrator/adaptor/jsonRpc"
)

type administrator struct {
	jsonRpcClient jsonRpcClient.Client
}

func New(
	jsonRpcClient jsonRpcClient.Client,
) administrator2.Administrator {
	return &administrator{
		jsonRpcClient: jsonRpcClient,
	}
}

func (a *administrator) Submit(request *administrator2.SubmitRequest) (*administrator2.SubmitResponse, error) {
	return nil, brainException.NotImplemented{}
}

func (a *administrator) ValidateFailTaskRequest(request *administrator2.FailTaskRequest) error {
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

func (a *administrator) FailTask(request *administrator2.FailTaskRequest) (*administrator2.FailTaskResponse, error) {
	if err := a.ValidateFailTaskRequest(request); err != nil {
		return nil, err
	}

	// create wrapped identifier
	wrappedZX303TaskIdentifier, err := wrappedIdentifier.Wrap(request.ZX303TaskIdentifier)
	if err != nil {
		return nil, brainException.Unexpected{Reasons: []string{"wrapping device identifier", err.Error()}}
	}

	// perform request
	failTaskResponse := jsonRpc.FailTaskResponse{}
	if err := a.jsonRpcClient.JsonRpcRequest(
		"ZX303TaskAdministrator.FailTask",
		jsonRpc.FailTaskRequest{
			ZX303TaskIdentifier: *wrappedZX303TaskIdentifier,
		},
		&failTaskResponse,
	); err != nil {
		return nil, brainException.Unexpected{Reasons: []string{"fail task error", err.Error()}}
	}

	return &administrator2.FailTaskResponse{
		ZX303Task: failTaskResponse.ZX303Task,
	}, nil
}

func (a *administrator) ValidateTransitionTaskRequest(request *administrator2.TransitionTaskRequest) error {
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

func (a *administrator) TransitionTask(request *administrator2.TransitionTaskRequest) (*administrator2.TransitionTaskResponse, error) {
	if err := a.ValidateTransitionTaskRequest(request); err != nil {
		return nil, err
	}

	// create wrapped identifier
	wrappedZX303TaskIdentifier, err := wrappedIdentifier.Wrap(request.ZX303TaskIdentifier)
	if err != nil {
		return nil, brainException.Unexpected{Reasons: []string{"wrapping device identifier", err.Error()}}
	}

	// perform request
	transitionTaskResponse := jsonRpc.TransitionTaskResponse{}
	if err := a.jsonRpcClient.JsonRpcRequest(
		"ZX303TaskAdministrator.TransitionTask",
		jsonRpc.TransitionTaskRequest{
			ZX303TaskIdentifier: *wrappedZX303TaskIdentifier,
			StepIdx:             request.StepIdx,
			NewStepStatus:       request.NewStepStatus,
		},
		&transitionTaskResponse,
	); err != nil {
		return nil, brainException.Unexpected{Reasons: []string{"fail task error", err.Error()}}
	}

	return &administrator2.TransitionTaskResponse{
		ZX303Task: transitionTaskResponse.ZX303Task,
	}, nil
}

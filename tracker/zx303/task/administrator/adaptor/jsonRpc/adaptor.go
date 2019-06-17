package jsonRpc

import (
	"github.com/iot-my-world/brain/log"
	wrappedIdentifier "github.com/iot-my-world/brain/search/identifier/wrapped"
	wrappedClaims "github.com/iot-my-world/brain/security/claims/wrapped"
	zx303Task "github.com/iot-my-world/brain/tracker/zx303/task"
	zx303TaskAdministrator "github.com/iot-my-world/brain/tracker/zx303/task/administrator"
	"github.com/iot-my-world/brain/tracker/zx303/task/step"
	"net/http"
)

type Adaptor struct {
	administrator zx303TaskAdministrator.Administrator
}

func New(administrator zx303TaskAdministrator.Administrator) *Adaptor {
	return &Adaptor{
		administrator: administrator,
	}
}

type SubmitRequest struct {
	ZX303Task zx303Task.Task `json:"zx303Task"`
}

type SubmitResponse struct {
	ZX303Task zx303Task.Task `json:"zx303Task"`
}

func (a *Adaptor) Submit(r *http.Request, request *SubmitRequest, response *SubmitResponse) error {
	claims, err := wrappedClaims.UnwrapClaimsFromContext(r)
	if err != nil {
		log.Warn(err.Error())
		return err
	}

	createResponse, err := a.administrator.Submit(&zx303TaskAdministrator.SubmitRequest{
		Claims:    claims,
		ZX303Task: request.ZX303Task,
	})
	if err != nil {
		return err
	}

	response.ZX303Task = createResponse.ZX303Task

	return nil
}

type FailTaskRequest struct {
	ZX303TaskIdentifier wrappedIdentifier.Wrapped `json:"zx303TaskIdentifier"`
}

type FailTaskResponse struct {
	ZX303Task zx303Task.Task `json:"zx303Task"`
}

func (a *Adaptor) FailTask(r *http.Request, request *FailTaskRequest, response *FailTaskResponse) error {
	claims, err := wrappedClaims.UnwrapClaimsFromContext(r)
	if err != nil {
		log.Warn(err.Error())
		return err
	}

	failTaskResponse, err := a.administrator.FailTask(&zx303TaskAdministrator.FailTaskRequest{
		Claims:              claims,
		ZX303TaskIdentifier: request.ZX303TaskIdentifier.Identifier,
	})
	if err != nil {
		return err
	}

	response.ZX303Task = failTaskResponse.ZX303Task

	return nil
}

type TransitionTaskRequest struct {
	ZX303TaskIdentifier wrappedIdentifier.Wrapped `json:"zx303TaskIdentifier"`
	StepIdx             int                       `json:"stepIdx"`
	NewStepStatus       step.Status               `json:"newStepStatus"`
}

type TransitionTaskResponse struct {
	ZX303Task zx303Task.Task `json:"zx303Task"`
}

func (a *Adaptor) TransitionTask(r *http.Request, request *TransitionTaskRequest, response *TransitionTaskResponse) error {
	claims, err := wrappedClaims.UnwrapClaimsFromContext(r)
	if err != nil {
		log.Warn(err.Error())
		return err
	}

	failTaskResponse, err := a.administrator.TransitionTask(&zx303TaskAdministrator.TransitionTaskRequest{
		Claims:              claims,
		ZX303TaskIdentifier: request.ZX303TaskIdentifier.Identifier,
		StepIdx:             request.StepIdx,
		NewStepStatus:       request.NewStepStatus,
	})
	if err != nil {
		return err
	}

	response.ZX303Task = failTaskResponse.ZX303Task

	return nil
}

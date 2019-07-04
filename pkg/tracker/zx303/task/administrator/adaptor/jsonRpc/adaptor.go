package jsonRpc

import (
	"github.com/iot-my-world/brain/internal/log"
	wrappedIdentifier "github.com/iot-my-world/brain/pkg/search/identifier/wrapped"
	"github.com/iot-my-world/brain/pkg/tracker/zx303/task"
	"github.com/iot-my-world/brain/pkg/tracker/zx303/task/administrator"
	step2 "github.com/iot-my-world/brain/pkg/tracker/zx303/task/step"
	wrappedClaims "github.com/iot-my-world/brain/security/claims/wrapped"
	"net/http"
)

type Adaptor struct {
	administrator administrator.Administrator
}

func New(administrator administrator.Administrator) *Adaptor {
	return &Adaptor{
		administrator: administrator,
	}
}

type SubmitRequest struct {
	ZX303Task task.Task `json:"zx303Task"`
}

type SubmitResponse struct {
	ZX303Task task.Task `json:"zx303Task"`
}

func (a *Adaptor) Submit(r *http.Request, request *SubmitRequest, response *SubmitResponse) error {
	claims, err := wrappedClaims.UnwrapClaimsFromContext(r)
	if err != nil {
		log.Warn(err.Error())
		return err
	}

	createResponse, err := a.administrator.Submit(&administrator.SubmitRequest{
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
	ZX303Task task.Task `json:"zx303Task"`
}

func (a *Adaptor) FailTask(r *http.Request, request *FailTaskRequest, response *FailTaskResponse) error {
	claims, err := wrappedClaims.UnwrapClaimsFromContext(r)
	if err != nil {
		log.Warn(err.Error())
		return err
	}

	failTaskResponse, err := a.administrator.FailTask(&administrator.FailTaskRequest{
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
	NewStepStatus       step2.Status              `json:"newStepStatus"`
}

type TransitionTaskResponse struct {
	ZX303Task task.Task `json:"zx303Task"`
}

func (a *Adaptor) TransitionTask(r *http.Request, request *TransitionTaskRequest, response *TransitionTaskResponse) error {
	claims, err := wrappedClaims.UnwrapClaimsFromContext(r)
	if err != nil {
		log.Warn(err.Error())
		return err
	}

	failTaskResponse, err := a.administrator.TransitionTask(&administrator.TransitionTaskRequest{
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

package administrator

import (
	"github.com/iot-my-world/brain/search/identifier"
	"github.com/iot-my-world/brain/security/claims"
	zx303Task "github.com/iot-my-world/brain/tracker/zx303/task"
	"github.com/iot-my-world/brain/tracker/zx303/task/step"
)

type Administrator interface {
	Submit(request *SubmitRequest) (*SubmitResponse, error)
	FailTask(request *FailTaskRequest) (*FailTaskResponse, error)
	TransitionTask(request *TransitionTaskRequest) (*TransitionTaskResponse, error)
}

type SubmitRequest struct {
	Claims    claims.Claims
	ZX303Task zx303Task.Task
}

type SubmitResponse struct {
	ZX303Task zx303Task.Task
}

type FailTaskRequest struct {
	Claims              claims.Claims
	ZX303TaskIdentifier identifier.Identifier
	FailedStepIdx       int
}

type FailTaskResponse struct {
	ZX303Task zx303Task.Task
}

type TransitionTaskRequest struct {
	Claims              claims.Claims
	ZX303TaskIdentifier identifier.Identifier
	StepIdx             int
	NewStepStatus       step.Status
}

type TransitionTaskResponse struct {
	ZX303Task zx303Task.Task
}

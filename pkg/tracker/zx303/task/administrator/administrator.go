package administrator

import (
	"github.com/iot-my-world/brain/pkg/search/identifier"
	"github.com/iot-my-world/brain/pkg/security/claims"
	"github.com/iot-my-world/brain/pkg/tracker/zx303/task"
	step2 "github.com/iot-my-world/brain/pkg/tracker/zx303/task/step"
)

type Administrator interface {
	Submit(request *SubmitRequest) (*SubmitResponse, error)
	FailTask(request *FailTaskRequest) (*FailTaskResponse, error)
	TransitionTask(request *TransitionTaskRequest) (*TransitionTaskResponse, error)
}

type SubmitRequest struct {
	Claims    claims.Claims
	ZX303Task task.Task
}

type SubmitResponse struct {
	ZX303Task task.Task
}

type FailTaskRequest struct {
	Claims              claims.Claims
	ZX303TaskIdentifier identifier.Identifier
	FailedStepIdx       int
}

type FailTaskResponse struct {
	ZX303Task task.Task
}

type TransitionTaskRequest struct {
	Claims              claims.Claims
	ZX303TaskIdentifier identifier.Identifier
	StepIdx             int
	NewStepStatus       step2.Status
}

type TransitionTaskResponse struct {
	ZX303Task task.Task
}

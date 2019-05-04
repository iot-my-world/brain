package administrator

import (
	"gitlab.com/iotTracker/brain/search/identifier"
	"gitlab.com/iotTracker/brain/security/claims"
	zx303Task "gitlab.com/iotTracker/brain/tracker/zx303/task"
	"gitlab.com/iotTracker/brain/tracker/zx303/task/step"
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

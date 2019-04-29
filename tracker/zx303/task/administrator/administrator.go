package administrator

import (
	"gitlab.com/iotTracker/brain/security/claims"
	zx303Task "gitlab.com/iotTracker/brain/tracker/zx303/task"
)

type Administrator interface {
	Submit(request *SubmitRequest) (*SubmitResponse, error)
}

type SubmitRequest struct {
	Claims    claims.Claims
	ZX303Task zx303Task.Task
}

type SubmitResponse struct {
	ZX303Task zx303Task.Task
}

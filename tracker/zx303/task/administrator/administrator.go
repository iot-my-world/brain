package administrator

import (
	"gitlab.com/iotTracker/brain/security/claims"
	zx303Task "gitlab.com/iotTracker/brain/tracker/zx303/task"
)

type Administrator interface {
	Create(request *CreateRequest) (*CreateResponse, error)
}

type CreateRequest struct {
	Claims    claims.Claims
	ZX303Task zx303Task.Task
}

type CreateResponse struct {
	ZX303Task zx303Task.Task
}

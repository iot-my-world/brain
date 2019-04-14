package administrator

import (
	"gitlab.com/iotTracker/brain/security/claims"
	"gitlab.com/iotTracker/brain/tracker/device/zx303"
)

type Administrator interface {
	Create(request *CreateRequest) (*CreateResponse, error)
}

type CreateRequest struct {
	Claims claims.Claims
	ZX303  zx303.ZX303
}

type CreateResponse struct {
	ZX303 zx303.ZX303
}

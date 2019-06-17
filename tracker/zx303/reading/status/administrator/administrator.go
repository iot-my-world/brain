package administrator

import (
	"gitlab.com/iotTracker/brain/security/claims"
	zx303StatusReading "gitlab.com/iotTracker/brain/tracker/zx303/reading/status"
)

type Administrator interface {
	Create(request *CreateRequest) (*CreateResponse, error)
}

type CreateRequest struct {
	Claims             claims.Claims
	ZX303StatusReading zx303StatusReading.Reading
}

type CreateResponse struct {
	ZX303StatusReading zx303StatusReading.Reading
}

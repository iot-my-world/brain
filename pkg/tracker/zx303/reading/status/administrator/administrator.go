package administrator

import (
	"github.com/iot-my-world/brain/pkg/tracker/zx303/reading/status"
	"github.com/iot-my-world/brain/security/claims"
)

type Administrator interface {
	Create(request *CreateRequest) (*CreateResponse, error)
}

type CreateRequest struct {
	Claims             claims.Claims
	ZX303StatusReading status.Reading
}

type CreateResponse struct {
	ZX303StatusReading status.Reading
}

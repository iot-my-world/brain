package administrator

import (
	"github.com/iot-my-world/brain/security/claims"
	zx303StatusReading "github.com/iot-my-world/brain/tracker/zx303/reading/status"
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

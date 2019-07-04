package administrator

import (
	"github.com/iot-my-world/brain/pkg/security/claims"
	"github.com/iot-my-world/brain/pkg/tracker/zx303/reading/gps"
)

type Administrator interface {
	Create(request *CreateRequest) (*CreateResponse, error)
}

type CreateRequest struct {
	Claims          claims.Claims
	ZX303GPSReading gps.Reading
}

type CreateResponse struct {
	ZX303GPSReading gps.Reading
}

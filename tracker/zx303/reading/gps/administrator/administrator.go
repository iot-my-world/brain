package administrator

import (
	"github.com/iot-my-world/brain/security/claims"
	zx303GPSReading "github.com/iot-my-world/brain/tracker/zx303/reading/gps"
)

type Administrator interface {
	Create(request *CreateRequest) (*CreateResponse, error)
}

type CreateRequest struct {
	Claims          claims.Claims
	ZX303GPSReading zx303GPSReading.Reading
}

type CreateResponse struct {
	ZX303GPSReading zx303GPSReading.Reading
}

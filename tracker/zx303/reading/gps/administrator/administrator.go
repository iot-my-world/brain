package administrator

import (
	"gitlab.com/iotTracker/brain/security/claims"
	zx303GPSReading "gitlab.com/iotTracker/brain/tracker/zx303/reading/gps"
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

package administrator

import (
	"gitlab.com/iotTracker/brain/security/claims"
	"gitlab.com/iotTracker/brain/tracker/device"
)

type Administrator interface {
	Create(request *CreateRequest) (*CreateResponse, error)
}

type CreateRequest struct {
	Claims claims.Claims
	Device device.Device
}

type CreateResponse struct {
	Device device.Device
}

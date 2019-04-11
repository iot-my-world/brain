package jsonRpc

import (
	"gitlab.com/iotTracker/brain/log"
	wrappedClaims "gitlab.com/iotTracker/brain/security/claims/wrapped"
	"gitlab.com/iotTracker/brain/tracker/device"
	wrappedDevice "gitlab.com/iotTracker/brain/tracker/device/wrapped"
	deviceAdministrator "gitlab.com/iotTracker/brain/tracker/device/administrator"
	"net/http"
)

type adaptor struct {
	deviceAdministrator deviceAdministrator.Administrator
}

func New(
	deviceAdministrator deviceAdministrator.Administrator,
) *adaptor {
	return &adaptor{
		deviceAdministrator: deviceAdministrator,
	}
}

type CreateRequest struct {
	WrappedDevice wrappedDevice.Wrapped `json:"device"`
}

type CreateResponse struct {
	Device device.Device `json:"device"`
}

func (a *adaptor) Create(r *http.Request, request *CreateRequest, response *CreateResponse) error {
	claims, err := wrappedClaims.UnwrapClaimsFromContext(r)
	if err != nil {
		log.Warn(err.Error())
		return err
	}

	createDeviceResponse, err := a.deviceAdministrator.Create(
		&deviceAdministrator.CreateRequest{
			Claims: claims,
			Device: request.WrappedDevice.Device,
		})
	if err != nil {
		return err
	}

	response.Device = createDeviceResponse.Device

	return nil
}

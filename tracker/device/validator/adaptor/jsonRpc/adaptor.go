package jsonRpc

import (
	"gitlab.com/iotTracker/brain/log"
	wrappedClaims "gitlab.com/iotTracker/brain/security/claims/wrapped"
	deviceValidator "gitlab.com/iotTracker/brain/tracker/device/validator"
	wrappedDevice "gitlab.com/iotTracker/brain/tracker/device/wrapped"
	"gitlab.com/iotTracker/brain/validate/reasonInvalid"
	"net/http"
)

type adaptor struct {
	deviceValidator deviceValidator.Validator
}

func New(
	deviceValidator deviceValidator.Validator,
) *adaptor {
	return &adaptor{
		deviceValidator: deviceValidator,
	}
}

type ValidateRequest struct {
	WrappedDevice wrappedDevice.Wrapped `json:"device"`
}

type ValidateResponse struct {
	ReasonsInvalid []reasonInvalid.ReasonInvalid `json:"reasonsInvalid"`
}

func (a *adaptor) Validate(r *http.Request, request *ValidateRequest, response *ValidateResponse) error {
	claims, err := wrappedClaims.UnwrapClaimsFromContext(r)
	if err != nil {
		log.Warn(err.Error())
		return err
	}

	createDeviceResponse, err := a.deviceValidator.Validate(
		&deviceValidator.ValidateRequest{
			Claims: claims,
			Device: request.WrappedDevice.Device,
		})
	if err != nil {
		return err
	}

	response.ReasonsInvalid = createDeviceResponse.ReasonsInvalid

	return nil
}

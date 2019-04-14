package basic

import (
	brainException "gitlab.com/iotTracker/brain/exception"
	"gitlab.com/iotTracker/brain/tracker/device"
	deviceValidator "gitlab.com/iotTracker/brain/tracker/device/validator"
)

type validator struct {
	deviceValidators map[device.Type]deviceValidator.Validator
}

func New(
	deviceValidators map[device.Type]deviceValidator.Validator,
) deviceValidator.Validator {
	return &validator{
		deviceValidators: deviceValidators,
	}
}

func (v *validator) ValidateValidateRequest(request *deviceValidator.ValidateRequest) error {
	reasonsInvalid := make([]string, 0)

	if request.Device == nil {
		reasonsInvalid = append(reasonsInvalid, "device is nil")
	} else {
		if v.deviceValidators[request.Device.Type()] == nil {
			reasonsInvalid = append(reasonsInvalid, "no validator for this device type")
		}
	}

	if request.Claims == nil {
		reasonsInvalid = append(reasonsInvalid, "claims are nil")
	}

	if len(reasonsInvalid) > 0 {
		return brainException.RequestInvalid{Reasons: reasonsInvalid}
	}
	return nil
}

func (v *validator) Validate(request *deviceValidator.ValidateRequest) (*deviceValidator.ValidateResponse, error) {
	if err := v.ValidateValidateRequest(request); err != nil {
		return nil, err
	}
	return v.deviceValidators[request.Device.Type()].Validate(request)
}

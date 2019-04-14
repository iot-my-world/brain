package basic

import (
	brainException "gitlab.com/iotTracker/brain/exception"
	deviceValidator "gitlab.com/iotTracker/brain/tracker/device/validator"
)

type validator struct {
}

func New() deviceValidator.Validator {
	return &validator{}
}

func (v *validator) ValidateValidateRequest(request *deviceValidator.ValidateRequest) error {
	reasonsInvalid := make([]string, 0)

	if len(reasonsInvalid) > 0 {
		return brainException.RequestInvalid{Reasons: reasonsInvalid}
	}
	return nil
}

func (v *validator) Validate(request *deviceValidator.ValidateRequest) (*deviceValidator.ValidateResponse, error) {
	if err := v.ValidateValidateRequest(request); err != nil {
		return nil, err
	}

	return &deviceValidator.ValidateResponse{}, nil
}

package basic

import (
	validator2 "github.com/iot-my-world/brain/pkg/tracker/tk102/reading/validator"
)

type validator struct {
}

func New() validator2.Validator {
	return &validator{}
}

func (v *validator) ValidateValidateRequest(request *validator2.ValidateRequest) error {
	return nil
}

func (v *validator) Validate(request *validator2.ValidateRequest) (*validator2.ValidateResponse, error) {
	if err := v.ValidateValidateRequest(request); err != nil {
		return nil, err
	}

	return &validator2.ValidateResponse{}, nil
}

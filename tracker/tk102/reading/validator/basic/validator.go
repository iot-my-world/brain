package basic

import (
	readingValidator "gitlab.com/iotTracker/brain/tracker/tk102/reading/validator"
)

type validator struct {
}

func New() readingValidator.Validator {
	return &validator{}
}

func (v *validator) ValidateValidateRequest(request *readingValidator.ValidateRequest) error {
	return nil
}

func (v *validator) Validate(request *readingValidator.ValidateRequest) (*readingValidator.ValidateResponse, error) {
	if err := v.ValidateValidateRequest(request); err != nil {
		return nil, err
	}

	return &readingValidator.ValidateResponse{}, nil
}

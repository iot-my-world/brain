package validator

import (
	"github.com/iot-my-world/brain/action"
	"github.com/iot-my-world/brain/party/individual"
	"github.com/iot-my-world/brain/security/claims"
	"github.com/iot-my-world/brain/validate/reasonInvalid"
)

type Validator interface {
	Validate(request *ValidateRequest) (*ValidateResponse, error)
}

type ValidateRequest struct {
	Claims     claims.Claims
	Individual individual.Individual
	Action     action.Action
}

type ValidateResponse struct {
	ReasonsInvalid []reasonInvalid.ReasonInvalid
}

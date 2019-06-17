package validator

import (
	"github.com/iot-my-world/brain/action"
	"github.com/iot-my-world/brain/party/client"
	"github.com/iot-my-world/brain/security/claims"
	"github.com/iot-my-world/brain/validate/reasonInvalid"
)

type Validator interface {
	Validate(request *ValidateRequest) (*ValidateResponse, error)
}

type ValidateRequest struct {
	Claims claims.Claims
	Client client.Client
	Action action.Action
}

type ValidateResponse struct {
	ReasonsInvalid []reasonInvalid.ReasonInvalid
}

package validator

import (
	"gitlab.com/iotTracker/brain/action"
	"gitlab.com/iotTracker/brain/party/client"
	"gitlab.com/iotTracker/brain/security/claims"
	"gitlab.com/iotTracker/brain/validate/reasonInvalid"
)

type Validator interface {
	Validate(request *ValidateRequest, response *ValidateResponse) error
}

type ValidateRequest struct {
	Claims claims.Claims
	Client client.Client
	Action action.Action
}

type ValidateResponse struct {
	ReasonsInvalid []reasonInvalid.ReasonInvalid
}

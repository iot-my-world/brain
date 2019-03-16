package validator

import (
	"gitlab.com/iotTracker/brain/action"
	"gitlab.com/iotTracker/brain/party/user"
	"gitlab.com/iotTracker/brain/security/claims"
	"gitlab.com/iotTracker/brain/validate/reasonInvalid"
)

type Validator interface {
	Validate(request *ValidateRequest, response *ValidateResponse) error
}

type ValidateRequest struct {
	Claims claims.Claims
	User   user.User
	Action action.Action
}

type ValidateResponse struct {
	ReasonsInvalid []reasonInvalid.ReasonInvalid
}

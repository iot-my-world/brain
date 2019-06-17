package validator

import (
	"gitlab.com/iotTracker/brain/action"
	"gitlab.com/iotTracker/brain/security/claims"
	humanUser "gitlab.com/iotTracker/brain/user/human"
	"gitlab.com/iotTracker/brain/validate/reasonInvalid"
)

type Validator interface {
	Validate(request *ValidateRequest) (*ValidateResponse, error)
}

type ValidateRequest struct {
	Claims claims.Claims
	User   humanUser.User
	Action action.Action
}

type ValidateResponse struct {
	ReasonsInvalid []reasonInvalid.ReasonInvalid
}

package validator

import (
	"gitlab.com/iotTracker/brain/action"
	"gitlab.com/iotTracker/brain/validate/reasonInvalid"
	"gitlab.com/iotTracker/brain/tracker/device/tk102"
	"gitlab.com/iotTracker/brain/security/claims"
)

type Validator interface {
	Validate(request *ValidateRequest, response *ValidateResponse) error
}

type ValidateRequest struct {
	Claims claims.Claims
	TK102  tk102.TK102
	Action action.Action
}

type ValidateResponse struct {
	ReasonsInvalid []reasonInvalid.ReasonInvalid
}

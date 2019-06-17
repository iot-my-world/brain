package validator

import (
	"gitlab.com/iotTracker/brain/action"
	"gitlab.com/iotTracker/brain/security/claims"
	zx303Task "gitlab.com/iotTracker/brain/tracker/zx303/task"
	"gitlab.com/iotTracker/brain/validate/reasonInvalid"
)

type Validator interface {
	Validate(request *ValidateRequest) (*ValidateResponse, error)
}

type ValidateRequest struct {
	Claims    claims.Claims
	ZX303Task zx303Task.Task
	Action    action.Action
}

type ValidateResponse struct {
	ReasonsInvalid []reasonInvalid.ReasonInvalid
}

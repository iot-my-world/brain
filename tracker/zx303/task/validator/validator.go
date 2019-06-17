package validator

import (
	"github.com/iot-my-world/brain/action"
	"github.com/iot-my-world/brain/security/claims"
	zx303Task "github.com/iot-my-world/brain/tracker/zx303/task"
	"github.com/iot-my-world/brain/validate/reasonInvalid"
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

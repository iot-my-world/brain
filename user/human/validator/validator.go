package validator

import (
	"github.com/iot-my-world/brain/action"
	"github.com/iot-my-world/brain/security/claims"
	humanUser "github.com/iot-my-world/brain/user/human"
	"github.com/iot-my-world/brain/validate/reasonInvalid"
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

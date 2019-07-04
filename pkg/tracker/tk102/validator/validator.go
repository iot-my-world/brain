package validator

import (
	"github.com/iot-my-world/brain/pkg/action"
	tk1022 "github.com/iot-my-world/brain/pkg/tracker/tk102"
	"github.com/iot-my-world/brain/security/claims"
	"github.com/iot-my-world/brain/validate/reasonInvalid"
)

type Validator interface {
	Validate(request *ValidateRequest) (*ValidateResponse, error)
}

type ValidateRequest struct {
	Claims claims.Claims
	TK102  tk1022.TK102
	Action action.Action
}

type ValidateResponse struct {
	ReasonsInvalid []reasonInvalid.ReasonInvalid
}

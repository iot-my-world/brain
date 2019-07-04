package validator

import (
	"github.com/iot-my-world/brain/pkg/action"
	zx3032 "github.com/iot-my-world/brain/pkg/tracker/zx303"
	"github.com/iot-my-world/brain/pkg/validate/reasonInvalid"
	"github.com/iot-my-world/brain/security/claims"
)

type Validator interface {
	Validate(request *ValidateRequest) (*ValidateResponse, error)
}

type ValidateRequest struct {
	Claims claims.Claims
	ZX303  zx3032.ZX303
	Action action.Action
}

type ValidateResponse struct {
	ReasonsInvalid []reasonInvalid.ReasonInvalid
}

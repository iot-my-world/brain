package validator

import (
	"github.com/iot-my-world/brain/pkg/action"
	"github.com/iot-my-world/brain/pkg/tracker/zx303/reading/status"
	"github.com/iot-my-world/brain/pkg/validate/reasonInvalid"
	"github.com/iot-my-world/brain/security/claims"
)

type Validator interface {
	Validate(request *ValidateRequest) (*ValidateResponse, error)
}

type ValidateRequest struct {
	Claims             claims.Claims
	ZX303StatusReading status.Reading
	Action             action.Action
}

type ValidateResponse struct {
	ReasonsInvalid []reasonInvalid.ReasonInvalid
}

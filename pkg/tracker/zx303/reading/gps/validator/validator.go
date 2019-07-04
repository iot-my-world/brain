package validator

import (
	"github.com/iot-my-world/brain/action"
	"github.com/iot-my-world/brain/pkg/tracker/zx303/reading/gps"
	"github.com/iot-my-world/brain/security/claims"
	"github.com/iot-my-world/brain/validate/reasonInvalid"
)

type Validator interface {
	Validate(request *ValidateRequest) (*ValidateResponse, error)
}

type ValidateRequest struct {
	Claims          claims.Claims
	ZX303GPSReading gps.Reading
	Action          action.Action
}

type ValidateResponse struct {
	ReasonsInvalid []reasonInvalid.ReasonInvalid
}

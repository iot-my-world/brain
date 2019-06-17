package validator

import (
	"github.com/iot-my-world/brain/action"
	"github.com/iot-my-world/brain/security/claims"
	zx303GPSReading "github.com/iot-my-world/brain/tracker/zx303/reading/gps"
	"github.com/iot-my-world/brain/validate/reasonInvalid"
)

type Validator interface {
	Validate(request *ValidateRequest) (*ValidateResponse, error)
}

type ValidateRequest struct {
	Claims          claims.Claims
	ZX303GPSReading zx303GPSReading.Reading
	Action          action.Action
}

type ValidateResponse struct {
	ReasonsInvalid []reasonInvalid.ReasonInvalid
}

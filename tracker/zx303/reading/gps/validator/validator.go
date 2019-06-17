package validator

import (
	"gitlab.com/iotTracker/brain/action"
	"gitlab.com/iotTracker/brain/security/claims"
	zx303GPSReading "gitlab.com/iotTracker/brain/tracker/zx303/reading/gps"
	"gitlab.com/iotTracker/brain/validate/reasonInvalid"
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

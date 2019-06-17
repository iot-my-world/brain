package validator

import (
	"github.com/iot-my-world/brain/action"
	"github.com/iot-my-world/brain/tracker/tk102/reading"
	"github.com/iot-my-world/brain/validate/reasonInvalid"
)

type Validator interface {
	Validate(request *ValidateRequest) (*ValidateResponse, error)
}

type ValidateRequest struct {
	Reading reading.Reading
	Action  action.Action
}

type ValidateResponse struct {
	ReasonsInvalid []reasonInvalid.ReasonInvalid
}

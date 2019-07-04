package validator

import (
	"github.com/iot-my-world/brain/pkg/action"
	reading2 "github.com/iot-my-world/brain/pkg/tracker/tk102/reading"
	"github.com/iot-my-world/brain/pkg/validate/reasonInvalid"
)

type Validator interface {
	Validate(request *ValidateRequest) (*ValidateResponse, error)
}

type ValidateRequest struct {
	Reading reading2.Reading
	Action  action.Action
}

type ValidateResponse struct {
	ReasonsInvalid []reasonInvalid.ReasonInvalid
}

package basic

import (
	"fmt"
	brainException "github.com/iot-my-world/brain/internal/exception"
	"github.com/iot-my-world/brain/pkg/tracker/zx303/reading/status/action"
	administrator2 "github.com/iot-my-world/brain/pkg/tracker/zx303/reading/status/administrator"
	"github.com/iot-my-world/brain/pkg/tracker/zx303/reading/status/administrator/exception"
	"github.com/iot-my-world/brain/pkg/tracker/zx303/reading/status/recordHandler"
	"github.com/iot-my-world/brain/pkg/tracker/zx303/reading/status/validator"
)

type administrator struct {
	zx303StatusReadingValidator     validator.Validator
	zx303StatusReadingRecordHandler *recordHandler.RecordHandler
}

func New(
	zx303StatusReadingValidator validator.Validator,
	zx303StatusReadingRecordHandler *recordHandler.RecordHandler,
) administrator2.Administrator {
	return &administrator{
		zx303StatusReadingValidator:     zx303StatusReadingValidator,
		zx303StatusReadingRecordHandler: zx303StatusReadingRecordHandler,
	}
}

func (a *administrator) ValidateCreateRequest(request *administrator2.CreateRequest) error {
	reasonsInvalid := make([]string, 0)

	if request.Claims == nil {
		reasonsInvalid = append(reasonsInvalid, "claims are nil")
	} else {
		zx303DeviceValidateResponse, err := a.zx303StatusReadingValidator.Validate(&validator.ValidateRequest{
			Claims:             request.Claims,
			ZX303StatusReading: request.ZX303StatusReading,
			Action:             action.Create,
		})
		if err != nil {
			reasonsInvalid = append(reasonsInvalid, "error validating zx303 status reading: "+err.Error())
		}
		if len(zx303DeviceValidateResponse.ReasonsInvalid) > 0 {
			for _, reason := range zx303DeviceValidateResponse.ReasonsInvalid {
				reasonsInvalid = append(reasonsInvalid, fmt.Sprintf("zx303 status reading invalid: %s - %s - %s", reason.Field, reason.Type, reason.Help))
			}
		}
	}

	if len(reasonsInvalid) > 0 {
		return brainException.RequestInvalid{Reasons: reasonsInvalid}
	}

	return nil
}

func (a *administrator) Create(request *administrator2.CreateRequest) (*administrator2.CreateResponse, error) {
	if err := a.ValidateCreateRequest(request); err != nil {
		return nil, err
	}

	createResponse, err := a.zx303StatusReadingRecordHandler.Create(&recordHandler.CreateRequest{
		ZX303StatusReading: request.ZX303StatusReading,
	})
	if err != nil {
		return nil, exception.ZX303StatusReadingCreation{Reasons: []string{err.Error()}}
	}

	return &administrator2.CreateResponse{
		ZX303StatusReading: createResponse.ZX303StatusReading,
	}, nil
}

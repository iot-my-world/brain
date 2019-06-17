package basic

import (
	"fmt"
	brainException "github.com/iot-my-world/brain/exception"
	zx303StatusReadingAction "github.com/iot-my-world/brain/tracker/zx303/reading/status/action"
	zx303StatusReadingAdministrator "github.com/iot-my-world/brain/tracker/zx303/reading/status/administrator"
	zx303StatusReadingAdministratorException "github.com/iot-my-world/brain/tracker/zx303/reading/status/administrator/exception"
	zx303StatusReadingRecordHandler "github.com/iot-my-world/brain/tracker/zx303/reading/status/recordHandler"
	zx303StatusReadingValidator "github.com/iot-my-world/brain/tracker/zx303/reading/status/validator"
)

type administrator struct {
	zx303StatusReadingValidator     zx303StatusReadingValidator.Validator
	zx303StatusReadingRecordHandler *zx303StatusReadingRecordHandler.RecordHandler
}

func New(
	zx303StatusReadingValidator zx303StatusReadingValidator.Validator,
	zx303StatusReadingRecordHandler *zx303StatusReadingRecordHandler.RecordHandler,
) zx303StatusReadingAdministrator.Administrator {
	return &administrator{
		zx303StatusReadingValidator:     zx303StatusReadingValidator,
		zx303StatusReadingRecordHandler: zx303StatusReadingRecordHandler,
	}
}

func (a *administrator) ValidateCreateRequest(request *zx303StatusReadingAdministrator.CreateRequest) error {
	reasonsInvalid := make([]string, 0)

	if request.Claims == nil {
		reasonsInvalid = append(reasonsInvalid, "claims are nil")
	} else {
		zx303DeviceValidateResponse, err := a.zx303StatusReadingValidator.Validate(&zx303StatusReadingValidator.ValidateRequest{
			Claims:             request.Claims,
			ZX303StatusReading: request.ZX303StatusReading,
			Action:             zx303StatusReadingAction.Create,
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

func (a *administrator) Create(request *zx303StatusReadingAdministrator.CreateRequest) (*zx303StatusReadingAdministrator.CreateResponse, error) {
	if err := a.ValidateCreateRequest(request); err != nil {
		return nil, err
	}

	createResponse, err := a.zx303StatusReadingRecordHandler.Create(&zx303StatusReadingRecordHandler.CreateRequest{
		ZX303StatusReading: request.ZX303StatusReading,
	})
	if err != nil {
		return nil, zx303StatusReadingAdministratorException.ZX303StatusReadingCreation{Reasons: []string{err.Error()}}
	}

	return &zx303StatusReadingAdministrator.CreateResponse{
		ZX303StatusReading: createResponse.ZX303StatusReading,
	}, nil
}

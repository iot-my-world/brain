package basic

import (
	"fmt"
	brainException "github.com/iot-my-world/brain/exception"
	"github.com/iot-my-world/brain/pkg/tracker/zx303/reading/gps/action"
	administrator2 "github.com/iot-my-world/brain/pkg/tracker/zx303/reading/gps/administrator"
	"github.com/iot-my-world/brain/pkg/tracker/zx303/reading/gps/administrator/exception"
	"github.com/iot-my-world/brain/pkg/tracker/zx303/reading/gps/recordHandler"
	"github.com/iot-my-world/brain/pkg/tracker/zx303/reading/gps/validator"
)

type administrator struct {
	zx303GPSReadingValidator     validator.Validator
	zx303GPSReadingRecordHandler *recordHandler.RecordHandler
}

func New(
	zx303GPSReadingValidator validator.Validator,
	zx303GPSReadingRecordHandler *recordHandler.RecordHandler,
) administrator2.Administrator {
	return &administrator{
		zx303GPSReadingValidator:     zx303GPSReadingValidator,
		zx303GPSReadingRecordHandler: zx303GPSReadingRecordHandler,
	}
}

func (a *administrator) ValidateCreateRequest(request *administrator2.CreateRequest) error {
	reasonsInvalid := make([]string, 0)

	if request.Claims == nil {
		reasonsInvalid = append(reasonsInvalid, "claims are nil")
	} else {
		zx303DeviceValidateResponse, err := a.zx303GPSReadingValidator.Validate(&validator.ValidateRequest{
			Claims:          request.Claims,
			ZX303GPSReading: request.ZX303GPSReading,
			Action:          action.Create,
		})
		if err != nil {
			reasonsInvalid = append(reasonsInvalid, "error validating zx303 gps reading: "+err.Error())
		}
		if len(zx303DeviceValidateResponse.ReasonsInvalid) > 0 {
			for _, reason := range zx303DeviceValidateResponse.ReasonsInvalid {
				reasonsInvalid = append(reasonsInvalid, fmt.Sprintf("zx303 gps reading invalid: %s - %s - %s", reason.Field, reason.Type, reason.Help))
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

	createResponse, err := a.zx303GPSReadingRecordHandler.Create(&recordHandler.CreateRequest{
		ZX303GPSReading: request.ZX303GPSReading,
	})
	if err != nil {
		return nil, exception.ZX303GPSReadingCreation{Reasons: []string{err.Error()}}
	}

	return &administrator2.CreateResponse{
		ZX303GPSReading: createResponse.ZX303GPSReading,
	}, nil
}

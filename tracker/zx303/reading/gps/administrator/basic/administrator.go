package basic

import (
	"fmt"
	brainException "gitlab.com/iotTracker/brain/exception"
	zx303GPSReadingAction "gitlab.com/iotTracker/brain/tracker/zx303/reading/gps/action"
	zx303GPSReadingAdministrator "gitlab.com/iotTracker/brain/tracker/zx303/reading/gps/administrator"
	zx303GPSReadingAdministratorException "gitlab.com/iotTracker/brain/tracker/zx303/reading/gps/administrator/exception"
	zx303GPSReadingRecordHandler "gitlab.com/iotTracker/brain/tracker/zx303/reading/gps/recordHandler"
	zx303GPSReadingValidator "gitlab.com/iotTracker/brain/tracker/zx303/reading/gps/validator"
)

type administrator struct {
	zx303GPSReadingValidator     zx303GPSReadingValidator.Validator
	zx303GPSReadingRecordHandler *zx303GPSReadingRecordHandler.RecordHandler
}

func New(
	zx303GPSReadingValidator zx303GPSReadingValidator.Validator,
	zx303GPSReadingRecordHandler *zx303GPSReadingRecordHandler.RecordHandler,
) zx303GPSReadingAdministrator.Administrator {
	return &administrator{
		zx303GPSReadingValidator:     zx303GPSReadingValidator,
		zx303GPSReadingRecordHandler: zx303GPSReadingRecordHandler,
	}
}

func (a *administrator) ValidateCreateRequest(request *zx303GPSReadingAdministrator.CreateRequest) error {
	reasonsInvalid := make([]string, 0)

	if request.Claims == nil {
		reasonsInvalid = append(reasonsInvalid, "claims are nil")
	} else {
		zx303DeviceValidateResponse, err := a.zx303GPSReadingValidator.Validate(&zx303GPSReadingValidator.ValidateRequest{
			Claims:          request.Claims,
			ZX303GPSReading: request.ZX303GPSReading,
			Action:          zx303GPSReadingAction.Create,
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

func (a *administrator) Create(request *zx303GPSReadingAdministrator.CreateRequest) (*zx303GPSReadingAdministrator.CreateResponse, error) {
	if err := a.ValidateCreateRequest(request); err != nil {
		return nil, err
	}

	createResponse, err := a.zx303GPSReadingRecordHandler.Create(&zx303GPSReadingRecordHandler.CreateRequest{
		ZX303GPSReading: request.ZX303GPSReading,
	})
	if err != nil {
		return nil, zx303GPSReadingAdministratorException.ZX303GPSReadingCreation{Reasons: []string{err.Error()}}
	}

	return &zx303GPSReadingAdministrator.CreateResponse{
		ZX303GPSReading: createResponse.ZX303GPSReading,
	}, nil
}

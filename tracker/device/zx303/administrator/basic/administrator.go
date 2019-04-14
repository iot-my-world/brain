package basic

import (
	"fmt"
	brainException "gitlab.com/iotTracker/brain/exception"
	zx303DeviceAction "gitlab.com/iotTracker/brain/tracker/device/zx303/action"
	zx303DeviceAdministrator "gitlab.com/iotTracker/brain/tracker/device/zx303/administrator"
	zx303DeviceAdministratorException "gitlab.com/iotTracker/brain/tracker/device/zx303/administrator/exception"
	zx303RecordHandler "gitlab.com/iotTracker/brain/tracker/device/zx303/recordHandler"
	zx303DeviceValidator "gitlab.com/iotTracker/brain/tracker/device/zx303/validator"
)

type administrator struct {
	zx303DeviceValidator zx303DeviceValidator.Validator
	zx303RecordHandler   *zx303RecordHandler.RecordHandler
}

func New(
	zx303DeviceValidator zx303DeviceValidator.Validator,
	zx303RecordHandler *zx303RecordHandler.RecordHandler,
) zx303DeviceAdministrator.Administrator {
	return &administrator{
		zx303DeviceValidator: zx303DeviceValidator,
		zx303RecordHandler:   zx303RecordHandler,
	}
}

func (a *administrator) ValidateCreateRequest(request *zx303DeviceAdministrator.CreateRequest) error {
	reasonsInvalid := make([]string, 0)

	if request.Claims == nil {
		reasonsInvalid = append(reasonsInvalid, "claims are nil")
	} else {
		zx303DeviceValidateResponse, err := a.zx303DeviceValidator.Validate(&zx303DeviceValidator.ValidateRequest{
			Claims: request.Claims,
			ZX303:  request.ZX303,
			Action: zx303DeviceAction.Create,
		})
		if err != nil {
			reasonsInvalid = append(reasonsInvalid, "error validating zx303 device: "+err.Error())
		}
		if len(zx303DeviceValidateResponse.ReasonsInvalid) > 0 {
			for _, reason := range zx303DeviceValidateResponse.ReasonsInvalid {
				reasonsInvalid = append(reasonsInvalid, fmt.Sprintf("zx303 device invalid: %s - %s - %s", reason.Field, reason.Type, reason.Help))
			}
		}
	}

	if len(reasonsInvalid) > 0 {
		return brainException.RequestInvalid{Reasons: reasonsInvalid}
	}

	return nil
}

func (a *administrator) Create(request *zx303DeviceAdministrator.CreateRequest) (*zx303DeviceAdministrator.CreateResponse, error) {
	if err := a.ValidateCreateRequest(request); err != nil {
		return nil, err
	}

	createResponse, err := a.zx303RecordHandler.Create(&zx303RecordHandler.CreateRequest{
		ZX303: request.ZX303,
	})
	if err != nil {
		return nil, zx303DeviceAdministratorException.DeviceCreation{Reasons: []string{err.Error()}}
	}

	return &zx303DeviceAdministrator.CreateResponse{
		ZX303: createResponse.ZX303,
	}, nil
}

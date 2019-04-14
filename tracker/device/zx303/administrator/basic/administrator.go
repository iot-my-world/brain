package basic

import (
	brainException "gitlab.com/iotTracker/brain/exception"
	zx303DeviceAdministrator "gitlab.com/iotTracker/brain/tracker/device/zx303/administrator"
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
		return nil, err
	}

	return &zx303DeviceAdministrator.CreateResponse{
		ZX303: createResponse.ZX303,
	}, nil
}

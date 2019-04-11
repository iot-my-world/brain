package basic

import (
	brainException "gitlab.com/iotTracker/brain/exception"
	deviceAdministrator "gitlab.com/iotTracker/brain/tracker/device/administrator"
	deviceRecordHandler "gitlab.com/iotTracker/brain/tracker/device/recordHandler"
)

type administrator struct {
	deviceRecordHandler *deviceRecordHandler.RecordHandler
}

func New(
	deviceRecordHandler *deviceRecordHandler.RecordHandler,
) deviceAdministrator.Administrator {
	return &administrator{
		deviceRecordHandler: deviceRecordHandler,
	}
}

func (a *administrator) ValidateCreateRequest(request *deviceAdministrator.CreateRequest) error {
	reasonsInvalid := make([]string, 0)

	if len(reasonsInvalid) > 0 {
		return brainException.RequestInvalid{Reasons: reasonsInvalid}
	}
	return nil
}

func (a *administrator) Create(request *deviceAdministrator.CreateRequest) (*deviceAdministrator.CreateResponse, error) {
	if err := a.ValidateCreateRequest(request); err != nil {
		return nil, err
	}

	return &deviceAdministrator.CreateResponse{}, nil
}

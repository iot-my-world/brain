package basic

import (
	brainException "gitlab.com/iotTracker/brain/exception"
	"gitlab.com/iotTracker/brain/tracker/device"
	deviceAdministrator "gitlab.com/iotTracker/brain/tracker/device/administrator"
	tk102DeviceRecordHandler "gitlab.com/iotTracker/brain/tracker/device/tk102/recordHandler"
)

type administrator struct {
	tk102DeviceRecordHandler tk102DeviceRecordHandler.RecordHandler
}

func New() deviceAdministrator.Administrator {
	return &administrator{}
}

func (a *administrator) ValidateCollectRequest(request *deviceAdministrator.CollectRequest) error {
	reasonsInvalid := make([]string, 0)

	if request.Claims == nil {
		reasonsInvalid = append(reasonsInvalid, "claims are nil")
	}

	if request.Criteria == nil {
		reasonsInvalid = append(reasonsInvalid, "criteria is nil")
	} else {
		for criterionIdx := range request.Criteria {
			if request.Criteria[criterionIdx] == nil {
				reasonsInvalid = append(reasonsInvalid, "a criterion is nil")
			}
		}
	}

	if len(reasonsInvalid) > 0 {
		return brainException.RequestInvalid{Reasons: reasonsInvalid}
	}

	return nil
}

func (a *administrator) Collect(request *deviceAdministrator.CollectRequest) (*deviceAdministrator.CollectResponse, error) {
	if err := a.ValidateCollectRequest(request); err != nil {
		return nil, err
	}

	return &deviceAdministrator.CollectResponse{
		Records: make([]device.Device, 0),
	}, nil
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

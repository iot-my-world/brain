package basic

import (
	clientRecordHandler "gitlab.com/iotTracker/brain/party/client/recordHandler"
	companyRecordHandler "gitlab.com/iotTracker/brain/party/company/recordHandler"
	tk102RecordHandler "gitlab.com/iotTracker/brain/tracker/device/tk102/recordHandler"
	tk102Administrator "gitlab.com/iotTracker/brain/tracker/device/tk102/administrator"
	brainException "gitlab.com/iotTracker/brain/exception"
)

type basicAdministrator struct {
	tk102RecordHandler   tk102RecordHandler.RecordHandler
	companyRecordHandler companyRecordHandler.RecordHandler
	clientRecordHandler  clientRecordHandler.RecordHandler
}

func New(
	tk102RecordHandler tk102RecordHandler.RecordHandler,
	companyRecordHandler companyRecordHandler.RecordHandler,
	clientRecordHandler clientRecordHandler.RecordHandler,
) *basicAdministrator {
	return &basicAdministrator{
		tk102RecordHandler:   tk102RecordHandler,
		companyRecordHandler: companyRecordHandler,
		clientRecordHandler:  clientRecordHandler,
	}
}

func (ba *basicAdministrator) ValidateChangeOwnerRequest(request *tk102Administrator.ChangeOwnerRequest) error {
	reasonsInvalid := make([]string, 0)

	if len(reasonsInvalid) > 0 {
		return brainException.RequestInvalid{Reasons: reasonsInvalid}
	} else {
		return nil
	}
}

func (ba *basicAdministrator) ChangeOwner(request *tk102Administrator.ChangeOwnerRequest, response *tk102Administrator.ChangeOwnerResponse) error {
	if err := ba.ValidateChangeOwnerRequest(request); err != nil {
		return err
	}

	return nil
}

func (ba *basicAdministrator) ValidateChangeAssignedRequest(request *tk102Administrator.ChangeAssignedRequest) error {
	reasonsInvalid := make([]string, 0)

	if len(reasonsInvalid) > 0 {
		return brainException.RequestInvalid{Reasons: reasonsInvalid}
	} else {
		return nil
	}
}

func (ba *basicAdministrator) ChangeAssigned(request *tk102Administrator.ChangeAssignedRequest, response *tk102Administrator.ChangeAssignedResponse) error {
	if err := ba.ValidateChangeAssignedRequest(request); err != nil {
		return err
	}

	return nil
}

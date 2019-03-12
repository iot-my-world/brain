package basic

import (
	brainException "gitlab.com/iotTracker/brain/exception"
	clientRecordHandler "gitlab.com/iotTracker/brain/party/client/recordHandler"
	companyRecordHandler "gitlab.com/iotTracker/brain/party/company/recordHandler"
	tk102Administrator "gitlab.com/iotTracker/brain/tracker/device/tk102/administrator"
	tk102RecordHandler "gitlab.com/iotTracker/brain/tracker/device/tk102/recordHandler"
)

type basicAdministrator struct {
	tk102RecordHandler   tk102RecordHandler.RecordHandler
	companyRecordHandler companyRecordHandler.RecordHandler
	clientRecordHandler  clientRecordHandler.RecordHandler
}

// New tk102 basic administrator
func New(
	tk102RecordHandler tk102RecordHandler.RecordHandler,
	companyRecordHandler companyRecordHandler.RecordHandler,
	clientRecordHandler clientRecordHandler.RecordHandler,
) tk102Administrator.Administrator {
	return &basicAdministrator{
		tk102RecordHandler:   tk102RecordHandler,
		companyRecordHandler: companyRecordHandler,
		clientRecordHandler:  clientRecordHandler,
	}
}

func (ba *basicAdministrator) ValidateChangeOwnerRequest(request *tk102Administrator.ChangeOwnerRequest) error {
	reasonsInvalid := make([]string, 0)

	if request.Claims == nil {
		reasonsInvalid = append(reasonsInvalid, "claims are nil")
	}
	if request.TK02Identifier == nil {
		reasonsInvalid = append(reasonsInvalid, "tk102 identifier is nil")
	}
	if request.NewOwnerIdentifier == nil {
		reasonsInvalid = append(reasonsInvalid, "newOwnerIdentifier is nil")
	}

	if len(reasonsInvalid) > 0 {
		return brainException.RequestInvalid{Reasons: reasonsInvalid}
	}
	return nil
}

func (ba *basicAdministrator) ChangeOwner(request *tk102Administrator.ChangeOwnerRequest, response *tk102Administrator.ChangeOwnerResponse) error {
	if err := ba.ValidateChangeOwnerRequest(request); err != nil {
		return err
	}

	// 1. retrieve the tk102 device

	return nil
}

func (ba *basicAdministrator) ValidateChangeAssignedRequest(request *tk102Administrator.ChangeAssignedRequest) error {
	reasonsInvalid := make([]string, 0)

	if len(reasonsInvalid) > 0 {
		return brainException.RequestInvalid{Reasons: reasonsInvalid}
	}
	return nil
}

func (ba *basicAdministrator) ChangeAssigned(request *tk102Administrator.ChangeAssignedRequest, response *tk102Administrator.ChangeAssignedResponse) error {
	if err := ba.ValidateChangeAssignedRequest(request); err != nil {
		return err
	}

	// Retrieve all readings associated with this device

	return nil
}

package basic

import (
	brainException "gitlab.com/iotTracker/brain/exception"
	clientRecordHandler "gitlab.com/iotTracker/brain/party/client/recordHandler"
	companyRecordHandler "gitlab.com/iotTracker/brain/party/company/recordHandler"
	partyHandler "gitlab.com/iotTracker/brain/party/handler"
	"gitlab.com/iotTracker/brain/search/criterion"
	exactTextCriterion "gitlab.com/iotTracker/brain/search/criterion/exact/text"
	tk102Administrator "gitlab.com/iotTracker/brain/tracker/device/tk102/administrator"
	tk102AdministratorException "gitlab.com/iotTracker/brain/tracker/device/tk102/administrator/exception"
	tk102RecordHandler "gitlab.com/iotTracker/brain/tracker/device/tk102/recordHandler"
	readingRecordHandler "gitlab.com/iotTracker/brain/tracker/reading/recordHandler"
)

type basicAdministrator struct {
	tk102RecordHandler   tk102RecordHandler.RecordHandler
	companyRecordHandler companyRecordHandler.RecordHandler
	clientRecordHandler  clientRecordHandler.RecordHandler
	partyHandler         partyHandler.Handler
	readingRecordHandler readingRecordHandler.RecordHandler
}

// New tk102 basic administrator
func New(
	tk102RecordHandler tk102RecordHandler.RecordHandler,
	companyRecordHandler companyRecordHandler.RecordHandler,
	clientRecordHandler clientRecordHandler.RecordHandler,
	partyHandler partyHandler.Handler,
	readingRecordHandler readingRecordHandler.RecordHandler,
) tk102Administrator.Administrator {
	return &basicAdministrator{
		tk102RecordHandler:   tk102RecordHandler,
		companyRecordHandler: companyRecordHandler,
		clientRecordHandler:  clientRecordHandler,
		partyHandler:         partyHandler,
		readingRecordHandler: readingRecordHandler,
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
	} else if request.Claims == nil {
		// confirm that the new owner exists by trying to retrieve their party
		if err := ba.partyHandler.RetrieveParty(&partyHandler.RetrievePartyRequest{
			Claims:     request.Claims,
			PartyType:  request.NewOwnerPartyType,
			Identifier: request.NewOwnerIdentifier,
		}, &partyHandler.RetrievePartyResponse{}); err != nil {
			reasonsInvalid = append(reasonsInvalid, "error retrieving new owners party: "+err.Error())
		}
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
	tk102RetrieveResponse := tk102RecordHandler.RetrieveResponse{}
	if err := ba.tk102RecordHandler.Retrieve(&tk102RecordHandler.RetrieveRequest{
		Claims:     request.Claims,
		Identifier: request.TK02Identifier,
	}, &tk102RetrieveResponse); err != nil {
		return tk102AdministratorException.DeviceRetrieval{Reasons: []string{err.Error()}}
	}

	// 2. collect readings for the device
	readingRetrieveResponse := readingRecordHandler.CollectResponse{}
	if err := ba.readingRecordHandler.Collect(&readingRecordHandler.CollectRequest{
		Claims: request.Claims,
		Criteria: []criterion.Criterion{
			exactTextCriterion.Criterion{
				Field: "id",
				Text:  tk102RetrieveResponse.TK102.Id,
			},
		},
		// Query: blank query as we have no restriction
	}, &readingRetrieveResponse); err != nil {
		return brainException.Unexpected{Reasons: []string{"collecting readings", err.Error()}}
	}

	// 3. update the device

	// 4. update the readings

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

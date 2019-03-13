package basic

import (
	brainException "gitlab.com/iotTracker/brain/exception"
	clientRecordHandler "gitlab.com/iotTracker/brain/party/client/recordHandler"
	companyRecordHandler "gitlab.com/iotTracker/brain/party/company/recordHandler"
	partyHandler "gitlab.com/iotTracker/brain/party/handler"
	"gitlab.com/iotTracker/brain/search/criterion"
	exactTextCriterion "gitlab.com/iotTracker/brain/search/criterion/exact/text"
	"gitlab.com/iotTracker/brain/search/identifier/id"
	tk102DeviceAdministrator "gitlab.com/iotTracker/brain/tracker/device/tk102/administrator"
	tk102DeviceAdministratorException "gitlab.com/iotTracker/brain/tracker/device/tk102/administrator/exception"
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
) tk102DeviceAdministrator.Administrator {
	return &basicAdministrator{
		tk102RecordHandler:   tk102RecordHandler,
		companyRecordHandler: companyRecordHandler,
		clientRecordHandler:  clientRecordHandler,
		partyHandler:         partyHandler,
		readingRecordHandler: readingRecordHandler,
	}
}

func (ba *basicAdministrator) ValidateChangeOwnerRequest(request *tk102DeviceAdministrator.ChangeOwnershipAndAssignmentRequest) error {
	reasonsInvalid := make([]string, 0)

	if request.Claims == nil {
		reasonsInvalid = append(reasonsInvalid, "claims are nil")
	}

	if len(reasonsInvalid) > 0 {
		return brainException.RequestInvalid{Reasons: reasonsInvalid}
	}
	return nil
}

/*
	1. retrieve the tk102 device
	2. collect readings for the device
	3. update the device
	4. update the readings
*/
func (ba *basicAdministrator) ChangeOwnershipAndAssignment(request *tk102DeviceAdministrator.ChangeOwnershipAndAssignmentRequest, response *tk102DeviceAdministrator.ChangeOwnershipAndAssignmentResponse) error {
	if err := ba.ValidateChangeOwnerRequest(request); err != nil {
		return err
	}

	// 1. retrieve the tk102 device
	tk102RetrieveResponse := tk102RecordHandler.RetrieveResponse{}
	if err := ba.tk102RecordHandler.Retrieve(&tk102RecordHandler.RetrieveRequest{
		Claims:     request.Claims,
		Identifier: id.Identifier{Id: request.TK02.Id},
	}, &tk102RetrieveResponse); err != nil {
		return tk102DeviceAdministratorException.DeviceRetrieval{Reasons: []string{err.Error()}}
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

	// 4. update the readings

	return nil
}

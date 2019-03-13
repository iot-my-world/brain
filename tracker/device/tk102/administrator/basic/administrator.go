package basic

import (
	"fmt"
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
	} else {
		// the device must be valid
		tk102ValidateResponse := tk102RecordHandler.ValidateResponse{}
		if err := ba.tk102RecordHandler.Validate(&tk102RecordHandler.ValidateRequest{
			Claims: request.Claims,
			TK102:  request.TK102,
		}, &tk102ValidateResponse); err != nil {
			reasonsInvalid = append(reasonsInvalid, "error validating device: "+err.Error())
		}
		if len(tk102ValidateResponse.ReasonsInvalid) > 0 {
			for _, reason := range tk102ValidateResponse.ReasonsInvalid {
				reasonsInvalid = append(reasonsInvalid, fmt.Sprintf("device invalid: %s - %s - %s", reason.Field, reason.Type, reason.Help))
			}
		} else {
			// if the assigned party is set OR
			if request.TK102.AssignedId.Id != "" ||
				// if the owner and assigned parties are not the same
				request.TK102.AssignedId.Id != request.TK102.OwnerId.Id {
				// then we must retrieve the owner and assigned parties to check the relationship is valid
				ownerPartyRetrieveResponse := partyHandler.RetrievePartyResponse{}
				if err := ba.partyHandler.RetrieveParty(&partyHandler.RetrievePartyRequest{
					Claims:     request.Claims,
					Identifier: request.TK102.OwnerId,
					PartyType:  request.TK102.OwnerPartyType,
				}, &ownerPartyRetrieveResponse); err != nil {
					reasonsInvalid = append(reasonsInvalid, "error retrieving owner party: "+err.Error())
				}
				assignedPartyRetrieveResponse := partyHandler.RetrievePartyResponse{}
				if err := ba.partyHandler.RetrieveParty(&partyHandler.RetrievePartyRequest{
					Claims:     request.Claims,
					Identifier: request.TK102.AssignedId,
					PartyType:  request.TK102.AssignedPartyType,
				}, &assignedPartyRetrieveResponse); err != nil {
					reasonsInvalid = append(reasonsInvalid, "error retrieving assigned party: "+err.Error())
				}

				// the owner party must be the parent party of the assigned party
				if ownerPartyRetrieveResponse.Party.Details().PartyType != assignedPartyRetrieveResponse.Party.Details().ParentPartyType ||
					ownerPartyRetrieveResponse.Party.Details().PartyId != assignedPartyRetrieveResponse.Party.Details().ParentId {
					reasonsInvalid = append(reasonsInvalid, "owner party must be the parent party of the assigned party")
				}
			}
		}
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
		Identifier: id.Identifier{Id: request.TK102.Id},
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
		return tk102DeviceAdministratorException.ReadingCollection{Reasons: []string{err.Error()}}
	}

	// 3. update the device
	tk102RetrieveResponse.TK102.OwnerPartyType = request.TK102.OwnerPartyType
	tk102RetrieveResponse.TK102.OwnerId = request.TK102.OwnerId
	tk102RetrieveResponse.TK102.AssignedPartyType = request.TK102.AssignedPartyType
	tk102RetrieveResponse.TK102.AssignedId = request.TK102.AssignedId
	tk102UpdateResponse := tk102RecordHandler.UpdateResponse{}
	if err := ba.tk102RecordHandler.Update(&tk102RecordHandler.UpdateRequest{
		Claims:     request.Claims,
		Identifier: id.Identifier{Id: request.TK102.Id},
		TK102:      tk102RetrieveResponse.TK102,
	}, &tk102UpdateResponse); err != nil {
		return tk102DeviceAdministratorException.DeviceUpdate{Reasons: []string{err.Error()}}
	}

	// 4. update the readings

	return nil
}

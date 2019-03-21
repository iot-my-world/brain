package basic

import (
	"fmt"
	brainException "gitlab.com/iotTracker/brain/exception"
	"gitlab.com/iotTracker/brain/party"
	partyAdministrator "gitlab.com/iotTracker/brain/party/administrator"
	clientRecordHandler "gitlab.com/iotTracker/brain/party/client/recordHandler"
	companyRecordHandler "gitlab.com/iotTracker/brain/party/company/recordHandler"
	"gitlab.com/iotTracker/brain/search/criterion"
	exactTextCriterion "gitlab.com/iotTracker/brain/search/criterion/exact/text"
	"gitlab.com/iotTracker/brain/search/identifier/id"
	tk102DeviceAction "gitlab.com/iotTracker/brain/tracker/device/tk102/action"
	tk102DeviceAdministrator "gitlab.com/iotTracker/brain/tracker/device/tk102/administrator"
	tk102DeviceAdministratorException "gitlab.com/iotTracker/brain/tracker/device/tk102/administrator/exception"
	tk102RecordHandler "gitlab.com/iotTracker/brain/tracker/device/tk102/recordHandler"
	tk102DeviceValidator "gitlab.com/iotTracker/brain/tracker/device/tk102/validator"
	readingRecordHandler "gitlab.com/iotTracker/brain/tracker/reading/recordHandler"
)

type administrator struct {
	tk102RecordHandler   tk102RecordHandler.RecordHandler
	companyRecordHandler companyRecordHandler.RecordHandler
	clientRecordHandler  clientRecordHandler.RecordHandler
	partyAdministrator   partyAdministrator.Administrator
	readingRecordHandler readingRecordHandler.RecordHandler
	tk102DeviceValidator tk102DeviceValidator.Validator
}

// New tk102 basic administrator
func New(
	tk102RecordHandler tk102RecordHandler.RecordHandler,
	companyRecordHandler companyRecordHandler.RecordHandler,
	clientRecordHandler clientRecordHandler.RecordHandler,
	partyAdministrator partyAdministrator.Administrator,
	readingRecordHandler readingRecordHandler.RecordHandler,
	tk102DeviceValidator tk102DeviceValidator.Validator,
) tk102DeviceAdministrator.Administrator {
	return &administrator{
		tk102RecordHandler:   tk102RecordHandler,
		companyRecordHandler: companyRecordHandler,
		clientRecordHandler:  clientRecordHandler,
		partyAdministrator:   partyAdministrator,
		readingRecordHandler: readingRecordHandler,
		tk102DeviceValidator: tk102DeviceValidator,
	}
}

func (a *administrator) ValidateChangeOwnershipAndAssignmentRequest(request *tk102DeviceAdministrator.ChangeOwnershipAndAssignmentRequest) error {
	reasonsInvalid := make([]string, 0)

	if request.Claims == nil {
		reasonsInvalid = append(reasonsInvalid, "claims are nil")
	} else {
		// the device must be valid
		tk102ValidateResponse, err := a.tk102DeviceValidator.Validate(&tk102DeviceValidator.ValidateRequest{
			Claims: request.Claims,
			TK102:  request.TK102,
			// Action: // no action. the device must be generally valid
		})
		if err != nil {
			reasonsInvalid = append(reasonsInvalid, "error validating device: "+err.Error())
		}
		if len(tk102ValidateResponse.ReasonsInvalid) > 0 {
			for _, reason := range tk102ValidateResponse.ReasonsInvalid {
				reasonsInvalid = append(reasonsInvalid, fmt.Sprintf("device invalid: %s - %s - %s", reason.Field, reason.Type, reason.Help))
			}
		} else {
			// if the party is not system, then the owner needs to be the party performing this request,
			// i.e. only system has the ability to change ownership
			if request.Claims.PartyDetails().PartyType != party.System {
				if !(request.TK102.OwnerId == request.Claims.PartyDetails().PartyId &&
					request.TK102.OwnerPartyType == request.Claims.PartyDetails().PartyType) {
					reasonsInvalid = append(reasonsInvalid, "only system can change tk102 device ownership")
				}
			}

			// if the assigned party is set OR
			if request.TK102.AssignedId.Id != "" ||
				// if the owner and assigned parties are not the same
				request.TK102.AssignedId.Id != request.TK102.OwnerId.Id {
				// then we must retrieve the owner and assigned parties to check the relationship is valid
				ownerPartyRetrieveResponse, err := a.partyAdministrator.RetrieveParty(&partyAdministrator.RetrievePartyRequest{
					Claims:     request.Claims,
					Identifier: request.TK102.OwnerId,
					PartyType:  request.TK102.OwnerPartyType,
				})
				if err != nil {
					reasonsInvalid = append(reasonsInvalid, "error retrieving owner party: "+err.Error())
				}
				assignedPartyRetrieveResponse, err := a.partyAdministrator.RetrieveParty(&partyAdministrator.RetrievePartyRequest{
					Claims:     request.Claims,
					Identifier: request.TK102.AssignedId,
					PartyType:  request.TK102.AssignedPartyType,
				})
				if err != nil {
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
ChangeOwnershipAndAssignment of a TK102 Tracking device
	1. retrieve the tk102 device
	2. collect readings for the device
	3. update the device
	4. update the readings
*/
func (a *administrator) ChangeOwnershipAndAssignment(request *tk102DeviceAdministrator.ChangeOwnershipAndAssignmentRequest) (*tk102DeviceAdministrator.ChangeOwnershipAndAssignmentResponse, error) {
	if err := a.ValidateChangeOwnershipAndAssignmentRequest(request); err != nil {
		return nil, err
	}

	// 1. retrieve the tk102 device
	tk102RetrieveResponse, err := a.tk102RecordHandler.Retrieve(&tk102RecordHandler.RetrieveRequest{
		Claims:     request.Claims,
		Identifier: id.Identifier{Id: request.TK102.Id},
	})
	if err != nil {
		return nil, tk102DeviceAdministratorException.DeviceRetrieval{Reasons: []string{err.Error()}}
	}

	// 2. collect readings for the device
	readingCollectResponse, err := a.readingRecordHandler.Collect(&readingRecordHandler.CollectRequest{
		Claims: request.Claims,
		Criteria: []criterion.Criterion{
			exactTextCriterion.Criterion{
				Field: "deviceId.id",
				Text:  tk102RetrieveResponse.TK102.Id,
			},
		},
		// Query: blank query as we have no restriction
	})
	if err != nil {
		return nil, tk102DeviceAdministratorException.ReadingCollection{Reasons: []string{err.Error()}}
	}

	// 3. update the device
	tk102RetrieveResponse.TK102.OwnerPartyType = request.TK102.OwnerPartyType
	tk102RetrieveResponse.TK102.OwnerId = request.TK102.OwnerId
	tk102RetrieveResponse.TK102.AssignedPartyType = request.TK102.AssignedPartyType
	tk102RetrieveResponse.TK102.AssignedId = request.TK102.AssignedId

	tk102UpdateResponse, err := a.tk102RecordHandler.Update(&tk102RecordHandler.UpdateRequest{
		Claims:     request.Claims,
		Identifier: id.Identifier{Id: request.TK102.Id},
		TK102:      tk102RetrieveResponse.TK102,
	})
	if err != nil {
		return nil, tk102DeviceAdministratorException.DeviceUpdate{Reasons: []string{err.Error()}}
	}

	// 4. update the readings
	for readingIdx := range readingCollectResponse.Records {
		readingCollectResponse.Records[readingIdx].OwnerPartyType = request.TK102.OwnerPartyType
		readingCollectResponse.Records[readingIdx].OwnerId = request.TK102.OwnerId
		readingCollectResponse.Records[readingIdx].AssignedPartyType = request.TK102.AssignedPartyType
		readingCollectResponse.Records[readingIdx].AssignedId = request.TK102.AssignedId
		if _, err := a.readingRecordHandler.Update(&readingRecordHandler.UpdateRequest{
			Claims:     request.Claims,
			Identifier: id.Identifier{Id: readingCollectResponse.Records[readingIdx].Id},
			Reading:    readingCollectResponse.Records[readingIdx],
		}); err != nil {
			return nil, tk102DeviceAdministratorException.ReadingUpdate{Reasons: []string{err.Error()}}
		}
	}

	return &tk102DeviceAdministrator.ChangeOwnershipAndAssignmentResponse{TK102: tk102UpdateResponse.TK102}, nil
}

func (a *administrator) ValidateCreateRequest(request *tk102DeviceAdministrator.CreateRequest) error {
	reasonsInvalid := make([]string, 0)

	if request.Claims == nil {
		reasonsInvalid = append(reasonsInvalid, "claims are nil")
	} else {
		tk102DeviceValidateResponse, err := a.tk102DeviceValidator.Validate(&tk102DeviceValidator.ValidateRequest{
			Claims: request.Claims,
			TK102:  request.TK102,
			Action: tk102DeviceAction.Create,
		})
		if err != nil {
			reasonsInvalid = append(reasonsInvalid, "error validating tk102 device: "+err.Error())
		}
		if len(tk102DeviceValidateResponse.ReasonsInvalid) > 0 {
			for _, reason := range tk102DeviceValidateResponse.ReasonsInvalid {
				reasonsInvalid = append(reasonsInvalid, fmt.Sprintf("tk102 device invalid: %s - %s - %s", reason.Field, reason.Type, reason.Help))
			}
		}
	}

	if len(reasonsInvalid) > 0 {
		return brainException.RequestInvalid{Reasons: reasonsInvalid}
	}
	return nil
}

func (a *administrator) Create(request *tk102DeviceAdministrator.CreateRequest) (*tk102DeviceAdministrator.CreateResponse, error) {
	if err := a.ValidateCreateRequest(request); err != nil {
		return nil, err
	}

	// Create the device
	createResponse, err := a.tk102RecordHandler.Create(&tk102RecordHandler.CreateRequest{
		TK102: request.TK102,
	})
	if err != nil {
		return nil, tk102DeviceAdministratorException.DeviceCreation{Reasons: []string{err.Error()}}
	}

	return &tk102DeviceAdministrator.CreateResponse{TK102: createResponse.TK102}, nil
}

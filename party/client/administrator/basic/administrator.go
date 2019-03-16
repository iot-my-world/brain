package basic

import (
	brainException "gitlab.com/iotTracker/brain/exception"
	"gitlab.com/iotTracker/brain/party"
	clientAdministrator "gitlab.com/iotTracker/brain/party/client/administrator"
	clientAdministratorException "gitlab.com/iotTracker/brain/party/client/administrator/exception"
	clientRecordHandler "gitlab.com/iotTracker/brain/party/client/recordHandler"
	"gitlab.com/iotTracker/brain/party/user"
	userRecordHandler "gitlab.com/iotTracker/brain/party/user/recordHandler"
	"gitlab.com/iotTracker/brain/search/identifier/id"
)

type administrator struct {
	clientRecordHandler clientRecordHandler.RecordHandler
	userRecordHandler   userRecordHandler.RecordHandler
}

func New(
	clientRecordHandler clientRecordHandler.RecordHandler,
	userRecordHandler userRecordHandler.RecordHandler,
) clientAdministrator.Administrator {
	return &administrator{
		clientRecordHandler: clientRecordHandler,
		userRecordHandler:   userRecordHandler,
	}
}

func (a *administrator) ValidateCreateResponse(request *clientAdministrator.CreateRequest) error {
	reasonsInvalid := make([]string, 0)

	if request.Claims == nil {
		reasonsInvalid = append(reasonsInvalid, "claims are nil")
	} else {
		switch request.Claims.PartyDetails().PartyType {
		case party.System:
			// do nothing, we expect system to not make a mistake
		default:
			if request.Client.ParentPartyType != request.Claims.PartyDetails().PartyType {
				reasonsInvalid = append(reasonsInvalid, "client ParentPartyType must be the type of the party doing creation")
			}
			if request.Client.ParentId != request.Claims.PartyDetails().PartyId {
				reasonsInvalid = append(reasonsInvalid, "client ParentId must be the id of the party doing creation")
			}
		}
	}

	if len(reasonsInvalid) > 0 {
		return brainException.RequestInvalid{Reasons: reasonsInvalid}
	}
	return nil
}

func (a *administrator) Create(request *clientAdministrator.CreateRequest, response *clientAdministrator.CreateResponse) error {
	if err := a.ValidateCreateResponse(request); err != nil {
		return err
	}

	// create minimal admin user for the client
	if err := a.userRecordHandler.Create(&userRecordHandler.CreateRequest{
		Claims: request.Claims,
		User: user.User{
			EmailAddress:    request.Client.AdminEmailAddress,
			ParentPartyType: request.Client.ParentPartyType,
			ParentId:        request.Client.ParentId,
			PartyType:       party.Client,
			PartyId:         id.Identifier{Id: request.Client.Id},
		},
	}, &userRecordHandler.CreateResponse{}); err != nil {
		return clientAdministratorException.ClientCreation{Reasons: []string{"creating admin user", err.Error()}}
	}

	// create the client itself

	return nil
}

func (a *administrator) ValidateUpdateAllowedFieldsRequest(request *clientAdministrator.UpdateAllowedFieldsRequest) error {
	reasonsInvalid := make([]string, 0)

	if request.Claims == nil {
		reasonsInvalid = append(reasonsInvalid, "claims are nil")
	}

	if len(reasonsInvalid) > 0 {
		return brainException.RequestInvalid{Reasons: reasonsInvalid}
	}
	return nil
}

func (a *administrator) UpdateAllowedFields(request *clientAdministrator.UpdateAllowedFieldsRequest, response *clientAdministrator.UpdateAllowedFieldsResponse) error {
	if err := a.ValidateUpdateAllowedFieldsRequest(request); err != nil {
		return err
	}

	// retrieve the client
	clientRetrieveResponse := clientRecordHandler.RetrieveResponse{}
	if err := a.clientRecordHandler.Retrieve(&clientRecordHandler.RetrieveRequest{
		Claims:     request.Claims,
		Identifier: id.Identifier{Id: request.Client.Id},
	}, &clientRetrieveResponse); err != nil {
		return clientAdministratorException.ClientRetrieval{Reasons: []string{err.Error()}}
	}

	// update the allowed fields on the client
	//clientRetrieveResponse.Client.Id = request.Client.Id
	//clientRetrieveResponse.Client.ParentId = request.Client.ParentId
	//clientRetrieveResponse.Client.ParentPartyType = request.Client.ParentPartyType
	clientRetrieveResponse.Client.Name = request.Client.Name
	//clientRetrieveResponse.Client.AdminEmailAddress = request.Client.AdminEmailAddress

	// update the client
	clientUpdateResponse := clientRecordHandler.UpdateResponse{}
	if err := a.clientRecordHandler.Update(&clientRecordHandler.UpdateRequest{
		Claims:     request.Claims,
		Identifier: id.Identifier{Id: request.Client.Id},
		Client:     clientRetrieveResponse.Client,
	}, &clientUpdateResponse); err != nil {
		return clientAdministratorException.AllowedFieldsUpdate{Reasons: []string{"updating", err.Error()}}
	}

	response.Client = clientUpdateResponse.Client

	return nil
}

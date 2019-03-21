package basic

import (
	"fmt"
	brainException "gitlab.com/iotTracker/brain/exception"
	"gitlab.com/iotTracker/brain/party"
	clientAction "gitlab.com/iotTracker/brain/party/client/action"
	clientAdministrator "gitlab.com/iotTracker/brain/party/client/administrator"
	clientAdministratorException "gitlab.com/iotTracker/brain/party/client/administrator/exception"
	clientRecordHandler "gitlab.com/iotTracker/brain/party/client/recordHandler"
	clientValidator "gitlab.com/iotTracker/brain/party/client/validator"
	"gitlab.com/iotTracker/brain/search/identifier/id"
	"gitlab.com/iotTracker/brain/user"
	userRecordHandler "gitlab.com/iotTracker/brain/user/recordHandler"
)

type administrator struct {
	clientRecordHandler clientRecordHandler.RecordHandler
	clientValidator     clientValidator.Validator
	userRecordHandler   userRecordHandler.RecordHandler
}

func New(
	clientRecordHandler clientRecordHandler.RecordHandler,
	clientValidator clientValidator.Validator,
	userRecordHandler userRecordHandler.RecordHandler,
) clientAdministrator.Administrator {
	return &administrator{
		clientRecordHandler: clientRecordHandler,
		clientValidator:     clientValidator,
		userRecordHandler:   userRecordHandler,
	}
}

func (a *administrator) ValidateCreateRequest(request *clientAdministrator.CreateRequest) error {
	reasonsInvalid := make([]string, 0)

	if request.Claims == nil {
		reasonsInvalid = append(reasonsInvalid, "claims are nil")
	} else {
		// confirm that the parent party of the client being created matches claims
		// i.e clients can only be created by their own parent party unless the system party
		// is acting
		switch request.Claims.PartyDetails().PartyType {
		case party.System:
			// do nothing, we expect system to know what they are doing
		default:
			if request.Client.ParentPartyType != request.Claims.PartyDetails().PartyType {
				reasonsInvalid = append(reasonsInvalid, "client ParentPartyType must be the type of the party doing creation")
			}
			if request.Client.ParentId != request.Claims.PartyDetails().PartyId {
				reasonsInvalid = append(reasonsInvalid, "client ParentId must be the id of the party doing creation")
			}
		}

		// client must be valid
		validationResponse := clientValidator.ValidateResponse{}
		if err := a.clientValidator.Validate(&clientValidator.ValidateRequest{
			Claims: request.Claims,
			Client: request.Client,
			Action: clientAction.Create,
		}, &validationResponse); err != nil {
			reasonsInvalid = append(reasonsInvalid, "error validating client: "+err.Error())
		}
		if len(validationResponse.ReasonsInvalid) > 0 {
			for _, reason := range validationResponse.ReasonsInvalid {
				reasonsInvalid = append(reasonsInvalid, fmt.Sprintf("client invalid: %s - %s - %s", reason.Field, reason.Type, reason.Help))
			}
		}
	}

	if len(reasonsInvalid) > 0 {
		return brainException.RequestInvalid{Reasons: reasonsInvalid}
	}
	return nil
}

func (a *administrator) Create(request *clientAdministrator.CreateRequest) (*clientAdministrator.CreateResponse, error) {
	if err := a.ValidateCreateRequest(request); err != nil {
		return nil, err
	}

	// create the client
	clientCreateResponse, err := a.clientRecordHandler.Create(&clientRecordHandler.CreateRequest{
		Client: request.Client,
	})
	if err != nil {
		return nil, clientAdministratorException.ClientCreation{Reasons: []string{"creating client", err.Error()}}
	}

	// create minimal admin user for the client
	if err := a.userRecordHandler.Create(&userRecordHandler.CreateRequest{
		User: user.User{
			EmailAddress:    clientCreateResponse.Client.AdminEmailAddress,
			ParentPartyType: clientCreateResponse.Client.ParentPartyType,
			ParentId:        clientCreateResponse.Client.ParentId,
			PartyType:       party.Client,
			PartyId:         id.Identifier{Id: clientCreateResponse.Client.Id},
		},
	}, &userRecordHandler.CreateResponse{}); err != nil {
		return nil, clientAdministratorException.ClientCreation{Reasons: []string{"creating admin user", err.Error()}}
	}

	return &clientAdministrator.CreateResponse{
		Client: clientCreateResponse.Client,
	}, nil
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

func (a *administrator) UpdateAllowedFields(request *clientAdministrator.UpdateAllowedFieldsRequest) (*clientAdministrator.UpdateAllowedFieldsResponse, error) {
	if err := a.ValidateUpdateAllowedFieldsRequest(request); err != nil {
		return nil, err
	}

	// retrieve the client
	clientRetrieveResponse, err := a.clientRecordHandler.Retrieve(&clientRecordHandler.RetrieveRequest{
		Claims:     request.Claims,
		Identifier: id.Identifier{Id: request.Client.Id},
	})
	if err != nil {
		return nil, clientAdministratorException.ClientRetrieval{Reasons: []string{err.Error()}}
	}

	// update the allowed fields on the client
	//clientRetrieveResponse.Client.Id = request.Client.Id
	//clientRetrieveResponse.Client.ParentId = request.Client.ParentId
	//clientRetrieveResponse.Client.ParentPartyType = request.Client.ParentPartyType
	clientRetrieveResponse.Client.Name = request.Client.Name
	//clientRetrieveResponse.Client.AdminEmailAddress = request.Client.AdminEmailAddress

	// update the client
	clientUpdateResponse, err := a.clientRecordHandler.Update(&clientRecordHandler.UpdateRequest{
		Claims:     request.Claims,
		Identifier: id.Identifier{Id: request.Client.Id},
		Client:     clientRetrieveResponse.Client,
	})
	if err != nil {
		return nil, clientAdministratorException.AllowedFieldsUpdate{Reasons: []string{"updating", err.Error()}}
	}

	return &clientAdministrator.UpdateAllowedFieldsResponse{
		Client: clientUpdateResponse.Client,
	}, nil
}

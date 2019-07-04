package basic

import (
	"fmt"
	brainException "github.com/iot-my-world/brain/exception"
	"github.com/iot-my-world/brain/log"
	"github.com/iot-my-world/brain/party"
	clientAction "github.com/iot-my-world/brain/party/client/action"
	clientAdministrator "github.com/iot-my-world/brain/party/client/administrator"
	clientAdministratorException "github.com/iot-my-world/brain/party/client/administrator/exception"
	clientRecordHandler "github.com/iot-my-world/brain/party/client/recordHandler"
	clientValidator "github.com/iot-my-world/brain/party/client/validator"
	"github.com/iot-my-world/brain/search/criterion"
	exactTextCriterion "github.com/iot-my-world/brain/search/criterion/exact/text"
	"github.com/iot-my-world/brain/search/identifier/id"
	humanUserLoginClaims "github.com/iot-my-world/brain/security/claims/login/user/human"
	humanUser "github.com/iot-my-world/brain/user/human"
	userRecordHandler "github.com/iot-my-world/brain/user/human/recordHandler"
)

type administrator struct {
	clientRecordHandler clientRecordHandler.RecordHandler
	clientValidator     clientValidator.Validator
	userRecordHandler   userRecordHandler.RecordHandler
	systemClaims        *humanUserLoginClaims.Login
}

func New(
	clientRecordHandler clientRecordHandler.RecordHandler,
	clientValidator clientValidator.Validator,
	userRecordHandler userRecordHandler.RecordHandler,
	systemClaims *humanUserLoginClaims.Login,
) clientAdministrator.Administrator {
	return &administrator{
		clientRecordHandler: clientRecordHandler,
		clientValidator:     clientValidator,
		userRecordHandler:   userRecordHandler,
		systemClaims:        systemClaims,
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
		validationResponse, err := a.clientValidator.Validate(&clientValidator.ValidateRequest{
			Claims: request.Claims,
			Client: request.Client,
			Action: clientAction.Create,
		})
		if err != nil {
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
	if _, err := a.userRecordHandler.Create(&userRecordHandler.CreateRequest{
		User: humanUser.User{
			EmailAddress:    clientCreateResponse.Client.AdminEmailAddress,
			ParentPartyType: clientCreateResponse.Client.ParentPartyType,
			ParentId:        clientCreateResponse.Client.ParentId,
			PartyType:       party.Client,
			PartyId:         id.Identifier{Id: clientCreateResponse.Client.Id},
		},
	}); err != nil {
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
	_, err = a.clientRecordHandler.Update(&clientRecordHandler.UpdateRequest{
		Claims:     request.Claims,
		Identifier: id.Identifier{Id: request.Client.Id},
		Client:     clientRetrieveResponse.Client,
	})
	if err != nil {
		return nil, clientAdministratorException.AllowedFieldsUpdate{Reasons: []string{"updating", err.Error()}}
	}

	return &clientAdministrator.UpdateAllowedFieldsResponse{
		Client: clientRetrieveResponse.Client,
	}, nil
}

func (a *administrator) ValidateDeleteRequest(request *clientAdministrator.DeleteRequest) error {
	reasonsInvalid := make([]string, 0)

	if request.ClientIdentifier == nil {
		reasonsInvalid = append(reasonsInvalid, "client identifier is nil")
	}

	if request.Claims == nil {
		reasonsInvalid = append(reasonsInvalid, "claims are nil")
	}

	if len(reasonsInvalid) > 0 {
		return brainException.RequestInvalid{Reasons: reasonsInvalid}
	}
	return nil
}

func (a *administrator) Delete(request *clientAdministrator.DeleteRequest) (*clientAdministrator.DeleteResponse, error) {
	if err := a.ValidateDeleteRequest(request); err != nil {
		log.Error(err.Error())
		return nil, err
	}

	// retrieve the client to be deleted
	clientRetrieveResponse, err := a.clientRecordHandler.Retrieve(&clientRecordHandler.RetrieveRequest{
		Claims:     request.Claims,
		Identifier: request.ClientIdentifier,
	})
	if err != nil {
		err = clientAdministratorException.Delete{Reasons: []string{"retrieve client error", err.Error()}}
		log.Error(err.Error())
		return nil, err
	}

	// collect any users in the client party
	clientUserCollectResponse, err := a.userRecordHandler.Collect(&userRecordHandler.CollectRequest{
		Claims: a.systemClaims, // using system claims since only system can see users from another party
		Criteria: []criterion.Criterion{
			exactTextCriterion.Criterion{
				Field: "partyId.id",
				Text:  clientRetrieveResponse.Client.Id,
			},
		},
	})
	if err != nil {
		err = clientAdministratorException.Delete{Reasons: []string{"collect users error", err.Error()}}
		log.Error(err.Error())
		return nil, err
	}

	// delete all users in the client party
	for idx := range clientUserCollectResponse.Records {
		if _, err := a.userRecordHandler.Delete(&userRecordHandler.DeleteRequest{
			Claims: a.systemClaims, // using system claims since only system can see users from another party
			Identifier: id.Identifier{
				Id: clientUserCollectResponse.Records[idx].Id,
			},
		}); err != nil {
			err = clientAdministratorException.Delete{Reasons: []string{"delete client user error", err.Error()}}
			log.Error(err.Error())
			return nil, err
		}
	}

	// delete client
	if _, err := a.clientRecordHandler.Delete(&clientRecordHandler.DeleteRequest{
		Claims:     request.Claims,
		Identifier: request.ClientIdentifier,
	}); err != nil {
		err = clientAdministratorException.Delete{Reasons: []string{"delete error", err.Error()}}
		log.Error(err.Error())
		return nil, err
	}

	return &clientAdministrator.DeleteResponse{}, nil
}

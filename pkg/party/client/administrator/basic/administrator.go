package basic

import (
	"fmt"
	brainException "github.com/iot-my-world/brain/internal/exception"
	"github.com/iot-my-world/brain/internal/log"
	"github.com/iot-my-world/brain/pkg/party"
	"github.com/iot-my-world/brain/pkg/party/client/action"
	administrator2 "github.com/iot-my-world/brain/pkg/party/client/administrator"
	"github.com/iot-my-world/brain/pkg/party/client/administrator/exception"
	"github.com/iot-my-world/brain/pkg/party/client/recordHandler"
	"github.com/iot-my-world/brain/pkg/party/client/validator"
	"github.com/iot-my-world/brain/pkg/search/criterion"
	exactTextCriterion "github.com/iot-my-world/brain/pkg/search/criterion/exact/text"
	"github.com/iot-my-world/brain/pkg/search/identifier/id"
	humanUserLoginClaims "github.com/iot-my-world/brain/pkg/security/claims/login/user/human"
	humanUser "github.com/iot-my-world/brain/pkg/user/human"
	userRecordHandler "github.com/iot-my-world/brain/pkg/user/human/recordHandler"
)

type administrator struct {
	clientRecordHandler recordHandler.RecordHandler
	clientValidator     validator.Validator
	userRecordHandler   userRecordHandler.RecordHandler
	systemClaims        *humanUserLoginClaims.Login
}

func New(
	clientRecordHandler recordHandler.RecordHandler,
	clientValidator validator.Validator,
	userRecordHandler userRecordHandler.RecordHandler,
	systemClaims *humanUserLoginClaims.Login,
) administrator2.Administrator {
	return &administrator{
		clientRecordHandler: clientRecordHandler,
		clientValidator:     clientValidator,
		userRecordHandler:   userRecordHandler,
		systemClaims:        systemClaims,
	}
}

func (a *administrator) ValidateCreateRequest(request *administrator2.CreateRequest) error {
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
		validationResponse, err := a.clientValidator.Validate(&validator.ValidateRequest{
			Claims: request.Claims,
			Client: request.Client,
			Action: action.Create,
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

func (a *administrator) Create(request *administrator2.CreateRequest) (*administrator2.CreateResponse, error) {
	if err := a.ValidateCreateRequest(request); err != nil {
		return nil, err
	}

	// create the client
	clientCreateResponse, err := a.clientRecordHandler.Create(&recordHandler.CreateRequest{
		Client: request.Client,
	})
	if err != nil {
		return nil, exception.ClientCreation{Reasons: []string{"creating client", err.Error()}}
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
		return nil, exception.ClientCreation{Reasons: []string{"creating admin user", err.Error()}}
	}

	return &administrator2.CreateResponse{
		Client: clientCreateResponse.Client,
	}, nil
}

func (a *administrator) ValidateUpdateAllowedFieldsRequest(request *administrator2.UpdateAllowedFieldsRequest) error {
	reasonsInvalid := make([]string, 0)

	if request.Claims == nil {
		reasonsInvalid = append(reasonsInvalid, "claims are nil")
	}

	if len(reasonsInvalid) > 0 {
		return brainException.RequestInvalid{Reasons: reasonsInvalid}
	}
	return nil
}

func (a *administrator) UpdateAllowedFields(request *administrator2.UpdateAllowedFieldsRequest) (*administrator2.UpdateAllowedFieldsResponse, error) {
	if err := a.ValidateUpdateAllowedFieldsRequest(request); err != nil {
		return nil, err
	}

	// retrieve the client
	clientRetrieveResponse, err := a.clientRecordHandler.Retrieve(&recordHandler.RetrieveRequest{
		Claims:     request.Claims,
		Identifier: id.Identifier{Id: request.Client.Id},
	})
	if err != nil {
		return nil, exception.ClientRetrieval{Reasons: []string{err.Error()}}
	}

	// update the allowed fields on the client
	//clientRetrieveResponse.Client.Id = request.Client.Id
	//clientRetrieveResponse.Client.ParentId = request.Client.ParentId
	//clientRetrieveResponse.Client.ParentPartyType = request.Client.ParentPartyType
	clientRetrieveResponse.Client.Name = request.Client.Name
	//clientRetrieveResponse.Client.AdminEmailAddress = request.Client.AdminEmailAddress

	// update the client
	_, err = a.clientRecordHandler.Update(&recordHandler.UpdateRequest{
		Claims:     request.Claims,
		Identifier: id.Identifier{Id: request.Client.Id},
		Client:     clientRetrieveResponse.Client,
	})
	if err != nil {
		return nil, exception.AllowedFieldsUpdate{Reasons: []string{"updating", err.Error()}}
	}

	return &administrator2.UpdateAllowedFieldsResponse{
		Client: clientRetrieveResponse.Client,
	}, nil
}

func (a *administrator) ValidateDeleteRequest(request *administrator2.DeleteRequest) error {
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

func (a *administrator) Delete(request *administrator2.DeleteRequest) (*administrator2.DeleteResponse, error) {
	if err := a.ValidateDeleteRequest(request); err != nil {
		log.Error(err.Error())
		return nil, err
	}

	// retrieve the client to be deleted
	clientRetrieveResponse, err := a.clientRecordHandler.Retrieve(&recordHandler.RetrieveRequest{
		Claims:     request.Claims,
		Identifier: request.ClientIdentifier,
	})
	if err != nil {
		err = exception.Delete{Reasons: []string{"retrieve client error", err.Error()}}
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
		err = exception.Delete{Reasons: []string{"collect users error", err.Error()}}
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
			err = exception.Delete{Reasons: []string{"delete client user error", err.Error()}}
			log.Error(err.Error())
			return nil, err
		}
	}

	// delete client
	if _, err := a.clientRecordHandler.Delete(&recordHandler.DeleteRequest{
		Claims:     request.Claims,
		Identifier: request.ClientIdentifier,
	}); err != nil {
		err = exception.Delete{Reasons: []string{"delete error", err.Error()}}
		log.Error(err.Error())
		return nil, err
	}

	return &administrator2.DeleteResponse{}, nil
}

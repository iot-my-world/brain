package basic

import (
	"fmt"
	brainException "github.com/iot-my-world/brain/exception"
	"github.com/iot-my-world/brain/search/identifier/id"
	apiUserAction "github.com/iot-my-world/brain/user/api/action"
	apiUserAdministrator "github.com/iot-my-world/brain/user/api/administrator"
	apiUserAdministratorException "github.com/iot-my-world/brain/user/api/administrator/exception"
	apiUserPasswordGenerator "github.com/iot-my-world/brain/user/api/password/generator"
	apiUserRecordHandler "github.com/iot-my-world/brain/user/api/recordHandler"
	apiUserValidator "github.com/iot-my-world/brain/user/api/validator"
	"github.com/satori/go.uuid"
	"golang.org/x/crypto/bcrypt"
)

type administrator struct {
	apiUserValidator         apiUserValidator.Validator
	apiUserRecordHandler     *apiUserRecordHandler.RecordHandler
	apiUserPasswordGenerator apiUserPasswordGenerator.Generator
}

func New(
	apiUserValidator apiUserValidator.Validator,
	apiUserRecordHandler *apiUserRecordHandler.RecordHandler,
	apiUserPasswordGenerator apiUserPasswordGenerator.Generator,
) apiUserAdministrator.Administrator {
	return &administrator{
		apiUserValidator:         apiUserValidator,
		apiUserRecordHandler:     apiUserRecordHandler,
		apiUserPasswordGenerator: apiUserPasswordGenerator,
	}
}

func (a *administrator) ValidateCreateRequest(request *apiUserAdministrator.CreateRequest) error {
	reasonsInvalid := make([]string, 0)

	if request.Claims == nil {
		reasonsInvalid = append(reasonsInvalid, "claims are nil")
	} else {
		apiUserAPIUserValidateResponse, err := a.apiUserValidator.Validate(&apiUserValidator.ValidateRequest{
			Claims: request.Claims,
			User:   request.User,
			Action: apiUserAction.Create,
		})
		if err != nil {
			reasonsInvalid = append(reasonsInvalid, "error validating apiUser apiUser: "+err.Error())
		}
		if len(apiUserAPIUserValidateResponse.ReasonsInvalid) > 0 {
			for _, reason := range apiUserAPIUserValidateResponse.ReasonsInvalid {
				reasonsInvalid = append(reasonsInvalid, fmt.Sprintf("apiUser apiUser invalid: %s - %s - %s", reason.Field, reason.Type, reason.Help))
			}
		}
	}

	if len(reasonsInvalid) > 0 {
		return brainException.RequestInvalid{Reasons: reasonsInvalid}
	}

	return nil
}

func (a *administrator) Create(request *apiUserAdministrator.CreateRequest) (*apiUserAdministrator.CreateResponse, error) {
	if err := a.ValidateCreateRequest(request); err != nil {
		return nil, err
	}

	// generate a username
	username, err := uuid.NewV4()
	if err != nil {
		return nil, brainException.UUIDGeneration{Reasons: []string{"username", err.Error()}}
	}

	// generate a password
	apiPasswordGenerateResponse, err := a.apiUserPasswordGenerator.Generate(&apiUserPasswordGenerator.GenerateRequest{
		CryptoBytesLength: 16,
	})
	if err != nil {
		return nil, apiUserAdministratorException.PasswordGeneration{Reasons: []string{err.Error()}}
	}

	// Hash the new Password
	pwdHash, err := bcrypt.GenerateFromPassword([]byte(apiPasswordGenerateResponse.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, apiUserAdministratorException.PasswordHash{Reasons: []string{err.Error()}}
	}

	request.User.Username = username.String()
	request.User.Password = pwdHash

	createResponse, err := a.apiUserRecordHandler.Create(&apiUserRecordHandler.CreateRequest{
		User: request.User,
	})
	if err != nil {
		return nil, apiUserAdministratorException.APIUserCreation{Reasons: []string{err.Error()}}
	}

	return &apiUserAdministrator.CreateResponse{
		User:     createResponse.User,
		Password: apiPasswordGenerateResponse.Password,
	}, nil
}

func (a *administrator) ValidateUpdateAllowedFieldsRequest(request *apiUserAdministrator.UpdateAllowedFieldsRequest) error {
	reasonsInvalid := make([]string, 0)

	if request.Claims == nil {
		reasonsInvalid = append(reasonsInvalid, "claims are nil")
	} else {
		// apiUser must be valid
		validationResponse, err := a.apiUserValidator.Validate(&apiUserValidator.ValidateRequest{
			Claims: request.Claims,
			Action: apiUserAction.UpdateAllowedFields,
		})
		if err != nil {
			reasonsInvalid = append(reasonsInvalid, "error validating apiUser: "+err.Error())
		}
		if len(validationResponse.ReasonsInvalid) > 0 {
			for _, reason := range validationResponse.ReasonsInvalid {
				reasonsInvalid = append(reasonsInvalid, fmt.Sprintf("apiUser apiUser invalid: %s - %s - %s", reason.Field, reason.Type, reason.Help))
			}
		}
	}

	if len(reasonsInvalid) > 0 {
		return brainException.RequestInvalid{Reasons: reasonsInvalid}
	}
	return nil
}

func (a *administrator) UpdateAllowedFields(request *apiUserAdministrator.UpdateAllowedFieldsRequest) (*apiUserAdministrator.UpdateAllowedFieldsResponse, error) {
	if err := a.ValidateUpdateAllowedFieldsRequest(request); err != nil {
		return nil, err
	}

	// retrieve the apiUser
	apiUserRetrieveResponse, err := a.apiUserRecordHandler.Retrieve(&apiUserRecordHandler.RetrieveRequest{
		Claims:     request.Claims,
		Identifier: id.Identifier{Id: request.User.Id},
	})
	if err != nil {
		return nil, apiUserAdministratorException.APIUserRetrieval{Reasons: []string{err.Error()}}
	}

	// update the allowed fields on the apiUser
	//apiUserRetrieveResponse.User.Id = request.User.Id
	apiUserRetrieveResponse.User.Name = request.User.Name
	apiUserRetrieveResponse.User.Description = request.User.Description
	//apiUserRetrieveResponse.User.Username = request.User.Username
	//apiUserRetrieveResponse.User.Password = request.User.Password
	//apiUserRetrieveResponse.User.Roles = request.User.Roles
	//apiUserRetrieveResponse.User.PartyType = request.User.PartyType
	//apiUserRetrieveResponse.User.PartyId = request.User.PartyId

	// update the apiUser
	_, err = a.apiUserRecordHandler.Update(&apiUserRecordHandler.UpdateRequest{
		Claims:     request.Claims,
		Identifier: id.Identifier{Id: request.User.Id},
		User:       apiUserRetrieveResponse.User,
	})
	if err != nil {
		return nil, apiUserAdministratorException.APIUserUpdate{Reasons: []string{err.Error()}}
	}

	return &apiUserAdministrator.UpdateAllowedFieldsResponse{
		User: apiUserRetrieveResponse.User,
	}, nil
}

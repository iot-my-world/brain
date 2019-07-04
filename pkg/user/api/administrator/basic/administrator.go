package basic

import (
	"fmt"
	brainException "github.com/iot-my-world/brain/exception"
	"github.com/iot-my-world/brain/pkg/search/identifier/id"
	"github.com/iot-my-world/brain/pkg/user/api/action"
	administrator2 "github.com/iot-my-world/brain/pkg/user/api/administrator"
	"github.com/iot-my-world/brain/pkg/user/api/administrator/exception"
	"github.com/iot-my-world/brain/pkg/user/api/password/generator"
	"github.com/iot-my-world/brain/pkg/user/api/recordHandler"
	"github.com/iot-my-world/brain/pkg/user/api/validator"
	"github.com/satori/go.uuid"
	"golang.org/x/crypto/bcrypt"
)

type administrator struct {
	apiUserValidator         validator.Validator
	apiUserRecordHandler     recordHandler.RecordHandler
	apiUserPasswordGenerator generator.Generator
}

func New(
	apiUserValidator validator.Validator,
	apiUserRecordHandler recordHandler.RecordHandler,
	apiUserPasswordGenerator generator.Generator,
) administrator2.Administrator {
	return &administrator{
		apiUserValidator:         apiUserValidator,
		apiUserRecordHandler:     apiUserRecordHandler,
		apiUserPasswordGenerator: apiUserPasswordGenerator,
	}
}

func (a *administrator) ValidateCreateRequest(request *administrator2.CreateRequest) error {
	reasonsInvalid := make([]string, 0)

	if request.Claims == nil {
		reasonsInvalid = append(reasonsInvalid, "claims are nil")
	} else {
		apiUserAPIUserValidateResponse, err := a.apiUserValidator.Validate(&validator.ValidateRequest{
			Claims: request.Claims,
			User:   request.User,
			Action: action.Create,
		})
		if err != nil {
			reasonsInvalid = append(reasonsInvalid, "error validating apiUser apiUser: "+err.Error())
		} else {
			if len(apiUserAPIUserValidateResponse.ReasonsInvalid) > 0 {
				for _, reason := range apiUserAPIUserValidateResponse.ReasonsInvalid {
					reasonsInvalid = append(reasonsInvalid, fmt.Sprintf("apiUser apiUser invalid: %s - %s - %s", reason.Field, reason.Type, reason.Help))
				}
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

	// generate a username
	username, err := uuid.NewV4()
	if err != nil {
		return nil, brainException.UUIDGeneration{Reasons: []string{"username", err.Error()}}
	}

	// generate a password
	apiPasswordGenerateResponse, err := a.apiUserPasswordGenerator.Generate(&generator.GenerateRequest{
		CryptoBytesLength: 16,
	})
	if err != nil {
		return nil, exception.PasswordGeneration{Reasons: []string{err.Error()}}
	}

	// Hash the new Password
	pwdHash, err := bcrypt.GenerateFromPassword([]byte(apiPasswordGenerateResponse.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, exception.PasswordHash{Reasons: []string{err.Error()}}
	}

	request.User.Username = username.String()
	request.User.Password = pwdHash

	createResponse, err := a.apiUserRecordHandler.Create(&recordHandler.CreateRequest{
		User: request.User,
	})
	if err != nil {
		return nil, exception.APIUserCreation{Reasons: []string{err.Error()}}
	}

	return &administrator2.CreateResponse{
		User:     createResponse.User,
		Password: apiPasswordGenerateResponse.Password,
	}, nil
}

func (a *administrator) ValidateUpdateAllowedFieldsRequest(request *administrator2.UpdateAllowedFieldsRequest) error {
	reasonsInvalid := make([]string, 0)

	if request.Claims == nil {
		reasonsInvalid = append(reasonsInvalid, "claims are nil")
	} else {
		// apiUser must be valid
		validationResponse, err := a.apiUserValidator.Validate(&validator.ValidateRequest{
			Claims: request.Claims,
			Action: action.UpdateAllowedFields,
		})
		if err != nil {
			reasonsInvalid = append(reasonsInvalid, "error validating apiUser: "+err.Error())
		} else {
			if len(validationResponse.ReasonsInvalid) > 0 {
				for _, reason := range validationResponse.ReasonsInvalid {
					reasonsInvalid = append(reasonsInvalid, fmt.Sprintf("apiUser apiUser invalid: %s - %s - %s", reason.Field, reason.Type, reason.Help))
				}
			}
		}
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

	// retrieve the apiUser
	apiUserRetrieveResponse, err := a.apiUserRecordHandler.Retrieve(&recordHandler.RetrieveRequest{
		Claims:     request.Claims,
		Identifier: id.Identifier{Id: request.User.Id},
	})
	if err != nil {
		return nil, exception.APIUserRetrieval{Reasons: []string{err.Error()}}
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
	_, err = a.apiUserRecordHandler.Update(&recordHandler.UpdateRequest{
		Claims:     request.Claims,
		Identifier: id.Identifier{Id: request.User.Id},
		User:       apiUserRetrieveResponse.User,
	})
	if err != nil {
		return nil, exception.APIUserUpdate{Reasons: []string{err.Error()}}
	}

	return &administrator2.UpdateAllowedFieldsResponse{
		User: apiUserRetrieveResponse.User,
	}, nil
}

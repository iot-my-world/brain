package basic

import (
	"fmt"
	brainException "gitlab.com/iotTracker/brain/exception"
	"gitlab.com/iotTracker/brain/search/identifier/id"
	apiUserAction "gitlab.com/iotTracker/brain/user/api/action"
	apiUserAdministrator "gitlab.com/iotTracker/brain/user/api/administrator"
	apiUserAdministratorException "gitlab.com/iotTracker/brain/user/api/administrator/exception"
	apiUserRecordHandler "gitlab.com/iotTracker/brain/user/api/recordHandler"
	apiUserValidator "gitlab.com/iotTracker/brain/user/api/validator"
)

type administrator struct {
	apiUserValidator     apiUserValidator.Validator
	apiUserRecordHandler *apiUserRecordHandler.RecordHandler
}

func New(
	apiUserValidator apiUserValidator.Validator,
	apiUserRecordHandler *apiUserRecordHandler.RecordHandler,
) apiUserAdministrator.Administrator {
	return &administrator{
		apiUserValidator:     apiUserValidator,
		apiUserRecordHandler: apiUserRecordHandler,
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

	createResponse, err := a.apiUserRecordHandler.Create(&apiUserRecordHandler.CreateRequest{
		User: request.User,
	})
	if err != nil {
		return nil, apiUserAdministratorException.APIUserCreation{Reasons: []string{err.Error()}}
	}

	return &apiUserAdministrator.CreateResponse{
		User: createResponse.User,
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

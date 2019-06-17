package basic

import (
	"fmt"
	brainException "github.com/iot-my-world/brain/exception"
	"github.com/iot-my-world/brain/party"
	companyAction "github.com/iot-my-world/brain/party/company/action"
	companyAdministrator "github.com/iot-my-world/brain/party/company/administrator"
	companyAdministratorException "github.com/iot-my-world/brain/party/company/administrator/exception"
	companyRecordHandler "github.com/iot-my-world/brain/party/company/recordHandler"
	companyValidator "github.com/iot-my-world/brain/party/company/validator"
	"github.com/iot-my-world/brain/search/identifier/id"
	humanUser "github.com/iot-my-world/brain/user/human"
	userRecordHandler "github.com/iot-my-world/brain/user/human/recordHandler"
)

type administrator struct {
	companyRecordHandler *companyRecordHandler.RecordHandler
	companyValidator     companyValidator.Validator
	userRecordHandler    userRecordHandler.RecordHandler
}

func New(
	companyRecordHandler *companyRecordHandler.RecordHandler,
	companyValidator companyValidator.Validator,
	userRecordHandler userRecordHandler.RecordHandler,
) companyAdministrator.Administrator {
	return &administrator{
		companyRecordHandler: companyRecordHandler,
		companyValidator:     companyValidator,
		userRecordHandler:    userRecordHandler,
	}
}

func (a *administrator) ValidateCreateRequest(request *companyAdministrator.CreateRequest) error {
	reasonsInvalid := make([]string, 0)

	// A new company can only be made by root
	if request.Claims == nil {
		reasonsInvalid = append(reasonsInvalid, "nil claims")
	} else {
		if request.Claims.PartyDetails().PartyType != party.System {
			reasonsInvalid = append(reasonsInvalid, "only system party can make a new company")
		}

		// company must be valid
		validationResponse, err := a.companyValidator.Validate(&companyValidator.ValidateRequest{
			Claims:  request.Claims,
			Company: request.Company,
			Action:  companyAction.Create,
		})
		if err != nil {
			reasonsInvalid = append(reasonsInvalid, "error validating company: "+err.Error())
		}
		if len(validationResponse.ReasonsInvalid) > 0 {
			for _, reason := range validationResponse.ReasonsInvalid {
				reasonsInvalid = append(reasonsInvalid, fmt.Sprintf("company invalid: %s - %s - %s", reason.Field, reason.Type, reason.Help))
			}
		}
	}

	if len(reasonsInvalid) > 0 {
		return brainException.RequestInvalid{Reasons: reasonsInvalid}
	}
	return nil
}

func (a *administrator) Create(request *companyAdministrator.CreateRequest) (*companyAdministrator.CreateResponse, error) {
	if err := a.ValidateCreateRequest(request); err != nil {
		return nil, err
	}

	// create the company
	companyCreateResponse, err := a.companyRecordHandler.Create(&companyRecordHandler.CreateRequest{
		Company: request.Company,
	})
	if err != nil {
		return nil, companyAdministratorException.CompanyCreation{Reasons: []string{"creating company", err.Error()}}
	}

	// create minimal admin user for the company
	if _, err := a.userRecordHandler.Create(&userRecordHandler.CreateRequest{
		User: humanUser.User{
			EmailAddress:    companyCreateResponse.Company.AdminEmailAddress,
			ParentPartyType: companyCreateResponse.Company.ParentPartyType,
			ParentId:        companyCreateResponse.Company.ParentId,
			PartyType:       party.Company,
			PartyId:         id.Identifier{Id: companyCreateResponse.Company.Id},
		},
	}); err != nil {
		return nil, companyAdministratorException.CompanyCreation{Reasons: []string{"creating admin user", err.Error()}}
	}

	return &companyAdministrator.CreateResponse{Company: companyCreateResponse.Company}, nil
}

func (a *administrator) ValidateUpdateAllowedFieldsRequest(request *companyAdministrator.UpdateAllowedFieldsRequest) error {
	reasonsInvalid := make([]string, 0)

	if request.Claims == nil {
		reasonsInvalid = append(reasonsInvalid, "claims are nil")
	}

	if len(reasonsInvalid) > 0 {
		return brainException.RequestInvalid{Reasons: reasonsInvalid}
	}
	return nil
}

func (a *administrator) UpdateAllowedFields(request *companyAdministrator.UpdateAllowedFieldsRequest) (*companyAdministrator.UpdateAllowedFieldsResponse, error) {
	if err := a.ValidateUpdateAllowedFieldsRequest(request); err != nil {
		return nil, err
	}

	// retrieve the company
	companyRetrieveResponse, err := a.companyRecordHandler.Retrieve(&companyRecordHandler.RetrieveRequest{
		Claims:     request.Claims,
		Identifier: id.Identifier{Id: request.Company.Id},
	})
	if err != nil {
		return nil, companyAdministratorException.CompanyRetrieval{Reasons: []string{err.Error()}}
	}

	// update the allowed fields on the company
	//companyRetrieveResponse.Company.Id = request.Company.Id
	//companyRetrieveResponse.Company.ParentId = request.Company.ParentId
	//companyRetrieveResponse.Company.ParentPartyType = request.Company.ParentPartyType
	companyRetrieveResponse.Company.Name = request.Company.Name
	//companyRetrieveResponse.Company.AdminEmailAddress = request.Company.AdminEmailAddress

	// update the company
	if _, err := a.companyRecordHandler.Update(&companyRecordHandler.UpdateRequest{
		Claims:     request.Claims,
		Identifier: id.Identifier{Id: request.Company.Id},
		Company:    companyRetrieveResponse.Company,
	}); err != nil {
		return nil, companyAdministratorException.AllowedFieldsUpdate{Reasons: []string{"updating", err.Error()}}
	}

	return &companyAdministrator.UpdateAllowedFieldsResponse{Company: companyRetrieveResponse.Company}, nil
}

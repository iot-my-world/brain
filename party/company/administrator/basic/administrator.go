package basic

import (
	brainException "gitlab.com/iotTracker/brain/exception"
	"gitlab.com/iotTracker/brain/party"
	companyAdministrator "gitlab.com/iotTracker/brain/party/company/administrator"
	companyAdministratorException "gitlab.com/iotTracker/brain/party/company/administrator/exception"
	companyRecordHandler "gitlab.com/iotTracker/brain/party/company/recordHandler"
	companyValidator "gitlab.com/iotTracker/brain/party/company/validator"
	"gitlab.com/iotTracker/brain/party/user"
	userRecordHandler "gitlab.com/iotTracker/brain/party/user/recordHandler"
	"gitlab.com/iotTracker/brain/search/identifier/id"
)

type administrator struct {
	companyRecordHandler companyRecordHandler.RecordHandler
	companyValidator     companyValidator.Validator
	userRecordHandler    userRecordHandler.RecordHandler
}

func New(
	companyRecordHandler companyRecordHandler.RecordHandler,
	companyValidator companyValidator.Validator,
	userRecordHandler userRecordHandler.RecordHandler,
) companyAdministrator.Administrator {
	return &administrator{
		companyRecordHandler: companyRecordHandler,
		companyValidator:     companyValidator,
		userRecordHandler:    userRecordHandler,
	}
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

func (a *administrator) ValidateCreateRequest(request *companyAdministrator.CreateRequest) error {
	reasonsInvalid := make([]string, 0)

	// A new company can only be made by root
	if request.Claims == nil {
		reasonsInvalid = append(reasonsInvalid, "nil claims")
	} else {
		if request.Claims.PartyDetails().PartyType != party.System {
			reasonsInvalid = append(reasonsInvalid, "only system party can make a new company")
		}
	}

	if len(reasonsInvalid) > 0 {
		return brainException.RequestInvalid{Reasons: reasonsInvalid}
	}
	return nil
}

func (a *administrator) Create(request *companyAdministrator.CreateRequest, response *companyAdministrator.CreateResponse) error {
	if err := a.ValidateCreateRequest(request); err != nil {
		return nil
	}

	// create minimal admin user for the company
	if err := a.userRecordHandler.Create(&userRecordHandler.CreateRequest{
		Claims: request.Claims,
		User: user.User{
			EmailAddress:    request.Company.AdminEmailAddress,
			ParentPartyType: request.Company.ParentPartyType,
			ParentId:        request.Company.ParentId,
			PartyType:       party.Company,
			PartyId:         id.Identifier{Id: request.Company.Id},
		},
	}, &userRecordHandler.CreateResponse{}); err != nil {
		return companyAdministratorException.CompanyCreation{Reasons: []string{"creating admin user", err.Error()}}
	}

	// create the company
	companyCreateResponse := companyRecordHandler.CreateResponse{}
	if err := a.companyRecordHandler.Create(&companyRecordHandler.CreateRequest{
		Claims:  request.Claims,
		Company: request.Company,
	}, &companyCreateResponse); err != nil {
		return companyAdministratorException.CompanyCreation{Reasons: []string{"creating company", err.Error()}}
	}

	response.Company = companyCreateResponse.Company

	return nil
}

func (a *administrator) UpdateAllowedFields(request *companyAdministrator.UpdateAllowedFieldsRequest, response *companyAdministrator.UpdateAllowedFieldsResponse) error {
	if err := a.ValidateUpdateAllowedFieldsRequest(request); err != nil {
		return err
	}

	// retrieve the company
	companyRetrieveResponse := companyRecordHandler.RetrieveResponse{}
	if err := a.companyRecordHandler.Retrieve(&companyRecordHandler.RetrieveRequest{
		Claims:     request.Claims,
		Identifier: id.Identifier{Id: request.Company.Id},
	}, &companyRetrieveResponse); err != nil {
		return companyAdministratorException.CompanyRetrieval{Reasons: []string{err.Error()}}
	}

	// update the allowed fields on the company
	//companyRetrieveResponse.Company.Id = request.Company.Id
	//companyRetrieveResponse.Company.ParentId = request.Company.ParentId
	//companyRetrieveResponse.Company.ParentPartyType = request.Company.ParentPartyType
	companyRetrieveResponse.Company.Name = request.Company.Name
	//companyRetrieveResponse.Company.AdminEmailAddress = request.Company.AdminEmailAddress

	// update the company
	companyUpdateResponse := companyRecordHandler.UpdateResponse{}
	if err := a.companyRecordHandler.Update(&companyRecordHandler.UpdateRequest{
		Claims:     request.Claims,
		Identifier: id.Identifier{Id: request.Company.Id},
		Company:    companyRetrieveResponse.Company,
	}, &companyUpdateResponse); err != nil {
		return companyAdministratorException.AllowedFieldsUpdate{Reasons: []string{"updating", err.Error()}}
	}

	response.Company = companyUpdateResponse.Company

	return nil
}

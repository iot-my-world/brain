package basic

import (
	brainException "gitlab.com/iotTracker/brain/exception"
	companyAdministrator "gitlab.com/iotTracker/brain/party/company/administrator"
	companyAdministratorException "gitlab.com/iotTracker/brain/party/company/administrator/exception"
	companyRecordHandler "gitlab.com/iotTracker/brain/party/company/recordHandler"
	"gitlab.com/iotTracker/brain/search/identifier/id"
)

type basicAdministrator struct {
	companyRecordHandler companyRecordHandler.RecordHandler
}

func New(
	companyRecordHandler companyRecordHandler.RecordHandler,
) companyAdministrator.Administrator {
	return &basicAdministrator{
		companyRecordHandler: companyRecordHandler,
	}
}

func (ba *basicAdministrator) ValidateUpdateAllowedFieldsRequest(request *companyAdministrator.UpdateAllowedFieldsRequest) error {
	reasonsInvalid := make([]string, 0)

	if request.Claims == nil {
		reasonsInvalid = append(reasonsInvalid, "claims are nil")
	}

	if len(reasonsInvalid) > 0 {
		return brainException.RequestInvalid{Reasons: reasonsInvalid}
	}
	return nil
}

func (ba *basicAdministrator) UpdateAllowedFields(request *companyAdministrator.UpdateAllowedFieldsRequest, response *companyAdministrator.UpdateAllowedFieldsResponse) error {
	if err := ba.ValidateUpdateAllowedFieldsRequest(request); err != nil {
		return err
	}

	// retrieve the company
	companyRetrieveResponse := companyRecordHandler.RetrieveResponse{}
	if err := ba.companyRecordHandler.Retrieve(&companyRecordHandler.RetrieveRequest{
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
	if err := ba.companyRecordHandler.Update(&companyRecordHandler.UpdateRequest{
		Claims:     request.Claims,
		Identifier: id.Identifier{Id: request.Company.Id},
		Company:    companyRetrieveResponse.Company,
	}, &companyUpdateResponse); err != nil {
		return companyAdministratorException.AllowedFieldsUpdate{Reasons: []string{"updating", err.Error()}}
	}

	response.Company = companyUpdateResponse.Company

	return nil
}

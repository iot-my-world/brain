package jsonRpc

import (
	"gitlab.com/iotTracker/brain/log"
	"gitlab.com/iotTracker/brain/party/company"
	companyAdministrator "gitlab.com/iotTracker/brain/party/company/administrator"
	"gitlab.com/iotTracker/brain/security/wrappedClaims"
	"net/http"
)

type adaptor struct {
	companyAdministrator companyAdministrator.Administrator
}

func New(
	companyAdministrator companyAdministrator.Administrator,
) *adaptor {
	return &adaptor{
		companyAdministrator: companyAdministrator,
	}
}

type CreateRequest struct {
	Company company.Company `json:"company"`
}

type CreateResponse struct {
	Company company.Company `json:"company"`
}

func (a *adaptor) Create(r *http.Request, request *CreateRequest, response *CreateResponse) error {
	claims, err := wrappedClaims.UnwrapClaimsFromContext(r)
	if err != nil {
		log.Warn(err.Error())
		return err
	}

	companyCreateResponse := companyAdministrator.CreateResponse{}
	if err := a.companyAdministrator.Create(&companyAdministrator.CreateRequest{
		Claims:  claims,
		Company: request.Company,
	}, &companyCreateResponse); err != nil {
		return err
	}

	response.Company = companyCreateResponse.Company

	return nil
}

type UpdateAllowedFieldsRequest struct {
	Company company.Company `json:"company"`
}

type UpdateAllowedFieldsResponse struct {
	Company company.Company `json:"company"`
}

func (a *adaptor) UpdateAllowedFields(r *http.Request, request *UpdateAllowedFieldsRequest, response *UpdateAllowedFieldsResponse) error {
	claims, err := wrappedClaims.UnwrapClaimsFromContext(r)
	if err != nil {
		log.Warn(err.Error())
		return err
	}

	updateAllowedFieldsResponse := companyAdministrator.UpdateAllowedFieldsResponse{}
	if err := a.companyAdministrator.UpdateAllowedFields(&companyAdministrator.UpdateAllowedFieldsRequest{
		Claims:  claims,
		Company: request.Company,
	}, &updateAllowedFieldsResponse); err != nil {
		return err
	}

	response.Company = updateAllowedFieldsResponse.Company

	return nil
}

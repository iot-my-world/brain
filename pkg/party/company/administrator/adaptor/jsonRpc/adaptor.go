package jsonRpc

import (
	"github.com/iot-my-world/brain/log"
	company2 "github.com/iot-my-world/brain/pkg/party/company"
	"github.com/iot-my-world/brain/pkg/party/company/administrator"
	wrappedIdentifier "github.com/iot-my-world/brain/search/identifier/wrapped"
	wrappedClaims "github.com/iot-my-world/brain/security/claims/wrapped"
	"net/http"
)

type adaptor struct {
	companyAdministrator administrator.Administrator
}

func New(
	companyAdministrator administrator.Administrator,
) *adaptor {
	return &adaptor{
		companyAdministrator: companyAdministrator,
	}
}

type CreateRequest struct {
	Company company2.Company `json:"company"`
}

type CreateResponse struct {
	Company company2.Company `json:"company"`
}

func (a *adaptor) Create(r *http.Request, request *CreateRequest, response *CreateResponse) error {
	claims, err := wrappedClaims.UnwrapClaimsFromContext(r)
	if err != nil {
		log.Warn(err.Error())
		return err
	}

	companyCreateResponse, err := a.companyAdministrator.Create(&administrator.CreateRequest{
		Claims:  claims,
		Company: request.Company,
	})
	if err != nil {
		return err
	}

	response.Company = companyCreateResponse.Company

	return nil
}

type UpdateAllowedFieldsRequest struct {
	Company company2.Company `json:"company"`
}

type UpdateAllowedFieldsResponse struct {
	Company company2.Company `json:"company"`
}

func (a *adaptor) UpdateAllowedFields(r *http.Request, request *UpdateAllowedFieldsRequest, response *UpdateAllowedFieldsResponse) error {
	claims, err := wrappedClaims.UnwrapClaimsFromContext(r)
	if err != nil {
		log.Warn(err.Error())
		return err
	}

	updateAllowedFieldsResponse, err := a.companyAdministrator.UpdateAllowedFields(&administrator.UpdateAllowedFieldsRequest{
		Claims:  claims,
		Company: request.Company,
	})
	if err != nil {
		return err
	}

	response.Company = updateAllowedFieldsResponse.Company

	return nil
}

type DeleteRequest struct {
	CompanyIdentifier wrappedIdentifier.Wrapped `json:"companyIdentifier"`
}

type DeleteResponse struct {
}

func (a *adaptor) Delete(r *http.Request, request *DeleteRequest, response *DeleteResponse) error {
	claims, err := wrappedClaims.UnwrapClaimsFromContext(r)
	if err != nil {
		log.Warn(err.Error())
		return err
	}

	if _, err := a.companyAdministrator.Delete(&administrator.DeleteRequest{
		Claims:            claims,
		CompanyIdentifier: request.CompanyIdentifier.Identifier,
	}); err != nil {
		return err
	}

	return nil
}

package administrator

import (
	"gitlab.com/iotTracker/brain/party/company"
	"gitlab.com/iotTracker/brain/security/claims"
)

type Administrator interface {
	UpdateAllowedFields(request *UpdateAllowedFieldsRequest, response *UpdateAllowedFieldsResponse) error
	Create(request *CreateRequest, response *CreateResponse) error
}

type UpdateAllowedFieldsRequest struct {
	Claims  claims.Claims
	Company company.Company
}

type UpdateAllowedFieldsResponse struct {
	Company company.Company
}

type CreateRequest struct {
	Claims  claims.Claims
	Company company.Company
}

type CreateResponse struct {
	Company company.Company
}

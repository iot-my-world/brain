package administrator

import (
	"github.com/iot-my-world/brain/party/company"
	"github.com/iot-my-world/brain/security/claims"
)

type Administrator interface {
	UpdateAllowedFields(request *UpdateAllowedFieldsRequest) (*UpdateAllowedFieldsResponse, error)
	Create(request *CreateRequest) (*CreateResponse, error)
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

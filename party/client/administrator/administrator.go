package administrator

import (
	"gitlab.com/iotTracker/brain/party/client"
	"gitlab.com/iotTracker/brain/security/claims"
)

type Administrator interface {
	UpdateAllowedFields(request *UpdateAllowedFieldsRequest) (*UpdateAllowedFieldsResponse, error)
	Create(request *CreateRequest) (*CreateResponse, error)
}

type CreateRequest struct {
	Claims claims.Claims
	Client client.Client
}

type CreateResponse struct {
	Client client.Client
}

type UpdateAllowedFieldsRequest struct {
	Claims claims.Claims
	Client client.Client
}

type UpdateAllowedFieldsResponse struct {
	Client client.Client
}

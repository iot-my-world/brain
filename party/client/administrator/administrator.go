package administrator

import (
	"gitlab.com/iotTracker/brain/party/client"
	"gitlab.com/iotTracker/brain/security/claims"
)

type Administrator interface {
	UpdateAllowedFields(request *UpdateAllowedFieldsRequest, response *UpdateAllowedFieldsResponse) error
}

type UpdateAllowedFieldsRequest struct {
	Claims claims.Claims
	Client client.Client
}

type UpdateAllowedFieldsResponse struct {
	Client client.Client
}

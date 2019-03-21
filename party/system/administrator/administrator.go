package administrator

import (
	"gitlab.com/iotTracker/brain/party/system"
	"gitlab.com/iotTracker/brain/security/claims"
)

type Administrator interface {
	UpdateAllowedFields(request *UpdateAllowedFieldsRequest) (*UpdateAllowedFieldsResponse, error)
}

type UpdateAllowedFieldsRequest struct {
	Claims claims.Claims
	System system.System
}

type UpdateAllowedFieldsResponse struct {
	System system.System
}

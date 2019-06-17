package administrator

import (
	"github.com/iot-my-world/brain/party/system"
	"github.com/iot-my-world/brain/security/claims"
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

package administrator

import (
	"gitlab.com/iotTracker/brain/search/identifier"
	"gitlab.com/iotTracker/brain/security/claims"
	"gitlab.com/iotTracker/brain/tracker/sf001"
)

type Administrator interface {
	Create(request *CreateRequest) (*CreateResponse, error)
	UpdateAllowedFields(request *UpdateAllowedFieldsRequest) (*UpdateAllowedFieldsResponse, error)
}

type CreateRequest struct {
	Claims claims.Claims
	SF001  sf001.SF001
}

type CreateResponse struct {
	SF001 sf001.SF001
}

type UpdateAllowedFieldsRequest struct {
	Claims claims.Claims
	SF001  sf001.SF001
}

type UpdateAllowedFieldsResponse struct {
	SF001 sf001.SF001
}

type HeartbeatRequest struct {
	Claims          claims.Claims
	SF001Identifier identifier.Identifier
}

type HeartbeatResponse struct {
}

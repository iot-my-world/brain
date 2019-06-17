package administrator

import (
	"gitlab.com/iotTracker/brain/search/identifier"
	"gitlab.com/iotTracker/brain/security/claims"
	"gitlab.com/iotTracker/brain/tracker/zx303"
)

type Administrator interface {
	Create(request *CreateRequest) (*CreateResponse, error)
	Heartbeat(request *HeartbeatRequest) (*HeartbeatResponse, error)
	UpdateAllowedFields(request *UpdateAllowedFieldsRequest) (*UpdateAllowedFieldsResponse, error)
}

type CreateRequest struct {
	Claims claims.Claims
	ZX303  zx303.ZX303
}

type CreateResponse struct {
	ZX303 zx303.ZX303
}

type UpdateAllowedFieldsRequest struct {
	Claims claims.Claims
	ZX303  zx303.ZX303
}

type UpdateAllowedFieldsResponse struct {
	ZX303 zx303.ZX303
}

type HeartbeatRequest struct {
	Claims          claims.Claims
	ZX303Identifier identifier.Identifier
}

type HeartbeatResponse struct {
}

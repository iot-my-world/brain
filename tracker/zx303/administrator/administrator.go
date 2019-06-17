package administrator

import (
	"github.com/iot-my-world/brain/search/identifier"
	"github.com/iot-my-world/brain/security/claims"
	"github.com/iot-my-world/brain/tracker/zx303"
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

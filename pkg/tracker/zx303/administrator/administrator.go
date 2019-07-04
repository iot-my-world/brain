package administrator

import (
	"github.com/iot-my-world/brain/pkg/search/identifier"
	"github.com/iot-my-world/brain/pkg/security/claims"
	zx3032 "github.com/iot-my-world/brain/pkg/tracker/zx303"
)

type Administrator interface {
	Create(request *CreateRequest) (*CreateResponse, error)
	Heartbeat(request *HeartbeatRequest) (*HeartbeatResponse, error)
	UpdateAllowedFields(request *UpdateAllowedFieldsRequest) (*UpdateAllowedFieldsResponse, error)
}

type CreateRequest struct {
	Claims claims.Claims
	ZX303  zx3032.ZX303
}

type CreateResponse struct {
	ZX303 zx3032.ZX303
}

type UpdateAllowedFieldsRequest struct {
	Claims claims.Claims
	ZX303  zx3032.ZX303
}

type UpdateAllowedFieldsResponse struct {
	ZX303 zx3032.ZX303
}

type HeartbeatRequest struct {
	Claims          claims.Claims
	ZX303Identifier identifier.Identifier
}

type HeartbeatResponse struct {
}

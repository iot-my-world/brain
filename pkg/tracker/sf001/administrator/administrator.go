package administrator

import (
	"github.com/iot-my-world/brain/pkg/search/identifier"
	sf0012 "github.com/iot-my-world/brain/pkg/tracker/sf001"
	"github.com/iot-my-world/brain/security/claims"
	"github.com/iot-my-world/brain/security/permission/api"
)

type Administrator interface {
	Create(request *CreateRequest) (*CreateResponse, error)
	UpdateAllowedFields(request *UpdateAllowedFieldsRequest) (*UpdateAllowedFieldsResponse, error)
}

const ServiceProvider = "SF001Tracker-Administrator"
const UpdateAllowedFieldsService = ServiceProvider + ".UpdateAllowedFields"
const CreateService = ServiceProvider + ".Create"

var SystemUserPermissions = []api.Permission{
	CreateService,
	UpdateAllowedFieldsService,
}

var CompanyAdminUserPermissions = make([]api.Permission, 0)

var CompanyUserPermissions = make([]api.Permission, 0)

var ClientAdminUserPermissions = make([]api.Permission, 0)

var ClientUserPermissions = make([]api.Permission, 0)

type CreateRequest struct {
	Claims claims.Claims
	SF001  sf0012.SF001
}

type CreateResponse struct {
	SF001 sf0012.SF001
}

type UpdateAllowedFieldsRequest struct {
	Claims claims.Claims
	SF001  sf0012.SF001
}

type UpdateAllowedFieldsResponse struct {
	SF001 sf0012.SF001
}

type HeartbeatRequest struct {
	Claims          claims.Claims
	SF001Identifier identifier.Identifier
}

type HeartbeatResponse struct {
}

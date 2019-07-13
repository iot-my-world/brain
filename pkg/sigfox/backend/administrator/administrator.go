package administrator

import (
	"github.com/iot-my-world/brain/pkg/search/identifier"
	"github.com/iot-my-world/brain/pkg/security/claims"
	"github.com/iot-my-world/brain/pkg/security/permission/api"
	"github.com/iot-my-world/brain/pkg/sigfox/backend"
)

type Administrator interface {
	Create(request *CreateRequest) (*CreateResponse, error)
	UpdateAllowedFields(request *UpdateAllowedFieldsRequest) (*UpdateAllowedFieldsResponse, error)
}

const ServiceProvider = "BackendTracker-Administrator"
const UpdateAllowedFieldsService = ServiceProvider + ".UpdateAllowedFields"
const CreateService = ServiceProvider + ".Create"

var SystemUserPermissions = []api.Permission{
	CreateService,
	UpdateAllowedFieldsService,
}

var CompanyAdminUserPermissions = []api.Permission{
	CreateService,
	UpdateAllowedFieldsService,
}

var CompanyUserPermissions = make([]api.Permission, 0)

var ClientAdminUserPermissions = []api.Permission{
	CreateService,
	UpdateAllowedFieldsService,
}

var ClientUserPermissions = make([]api.Permission, 0)

type CreateRequest struct {
	Claims  claims.Claims
	Backend backend.Backend
}

type CreateResponse struct {
	Backend backend.Backend
}

type UpdateAllowedFieldsRequest struct {
	Claims  claims.Claims
	Backend backend.Backend
}

type UpdateAllowedFieldsResponse struct {
	Backend backend.Backend
}

type HeartbeatRequest struct {
	Claims            claims.Claims
	BackendIdentifier identifier.Identifier
}

type HeartbeatResponse struct {
}

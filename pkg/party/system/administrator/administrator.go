package administrator

import (
	system2 "github.com/iot-my-world/brain/pkg/party/system"
	"github.com/iot-my-world/brain/pkg/security/claims"
	"github.com/iot-my-world/brain/pkg/security/permission/api"
)

type Administrator interface {
	UpdateAllowedFields(request *UpdateAllowedFieldsRequest) (*UpdateAllowedFieldsResponse, error)
}

const ServiceProvider = "System-Administrator"
const UpdateAllowedFieldsService = ServiceProvider + ".UpdateAllowedFields"

var SystemUserPermissions = []api.Permission{
	UpdateAllowedFieldsService,
}

var CompanyAdminUserPermissions = make([]api.Permission, 0)

var CompanyUserPermissions = make([]api.Permission, 0)

var ClientAdminUserPermissions = make([]api.Permission, 0)

var ClientUserPermissions = make([]api.Permission, 0)

type UpdateAllowedFieldsRequest struct {
	Claims claims.Claims
	System system2.System
}

type UpdateAllowedFieldsResponse struct {
	System system2.System
}

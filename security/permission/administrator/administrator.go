package administrator

import (
	"github.com/iot-my-world/brain/pkg/search/identifier"
	"github.com/iot-my-world/brain/security/claims"
	"github.com/iot-my-world/brain/security/permission/api"
	"github.com/iot-my-world/brain/security/permission/view"
)

type Administrator interface {
	UserHasPermission(request *UserHasPermissionRequest) (*UserHasPermissionResponse, error)
	GetAllUsersAPIPermissions(request *GetAllUsersAPIPermissionsRequest) (*GetAllUsersAPIPermissionsResponse, error)
	GetAllUsersViewPermissions(request *GetAllUsersViewPermissionsRequest) (*GetAllUsersViewPermissionsResponse, error)
}

const ServiceProvider = "Permission-Administrator"
const UserHasPermissionService = ServiceProvider + ".UserHasPermission"
const GetAllUsersAPIPermissionsService = ServiceProvider + ".GetAllUsersAPIPermissions"
const GetAllUsersViewPermissionsService = ServiceProvider + ".GetAllUsersViewPermissions"

var SystemUserPermissions = make([]api.Permission, 0)

var CompanyAdminUserPermissions = []api.Permission{
	GetAllUsersViewPermissionsService,
}

var CompanyUserPermissions = []api.Permission{
	GetAllUsersViewPermissionsService,
}

var ClientAdminUserPermissions = []api.Permission{
	GetAllUsersViewPermissionsService,
}

var ClientUserPermissions = []api.Permission{
	GetAllUsersViewPermissionsService,
}

type UserHasPermissionRequest struct {
	Claims         claims.Claims
	UserIdentifier identifier.Identifier
	Permission     api.Permission
}

type UserHasPermissionResponse struct {
	Result bool
}

type GetAllUsersAPIPermissionsRequest struct {
	Claims         claims.Claims
	UserIdentifier identifier.Identifier
}

type GetAllUsersAPIPermissionsResponse struct {
	Permissions []api.Permission
}

type GetAllUsersViewPermissionsRequest struct {
	Claims         claims.Claims
	UserIdentifier identifier.Identifier
}

type GetAllUsersViewPermissionsResponse struct {
	Permissions []view.Permission
}

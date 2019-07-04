package administrator

import (
	"github.com/iot-my-world/brain/pkg/search/identifier"
	claims2 "github.com/iot-my-world/brain/pkg/security/claims"
	api2 "github.com/iot-my-world/brain/pkg/security/permission/api"
	view2 "github.com/iot-my-world/brain/pkg/security/permission/view"
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

var SystemUserPermissions = make([]api2.Permission, 0)

var CompanyAdminUserPermissions = []api2.Permission{
	GetAllUsersViewPermissionsService,
}

var CompanyUserPermissions = []api2.Permission{
	GetAllUsersViewPermissionsService,
}

var ClientAdminUserPermissions = []api2.Permission{
	GetAllUsersViewPermissionsService,
}

var ClientUserPermissions = []api2.Permission{
	GetAllUsersViewPermissionsService,
}

type UserHasPermissionRequest struct {
	Claims         claims2.Claims
	UserIdentifier identifier.Identifier
	Permission     api2.Permission
}

type UserHasPermissionResponse struct {
	Result bool
}

type GetAllUsersAPIPermissionsRequest struct {
	Claims         claims2.Claims
	UserIdentifier identifier.Identifier
}

type GetAllUsersAPIPermissionsResponse struct {
	Permissions []api2.Permission
}

type GetAllUsersViewPermissionsRequest struct {
	Claims         claims2.Claims
	UserIdentifier identifier.Identifier
}

type GetAllUsersViewPermissionsResponse struct {
	Permissions []view2.Permission
}

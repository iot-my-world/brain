package administrator

import (
	"github.com/iot-my-world/brain/search/identifier"
	"github.com/iot-my-world/brain/security/claims"
	"github.com/iot-my-world/brain/security/permission/api"
	"github.com/iot-my-world/brain/security/permission/view"
)

type Administrator interface {
	UserHasPermission(request *UserHasPermissionRequest) (*UserHasPermissionResponse, error)
	GetAllUsersAPIPermissions(request *GetAllUsersAPIPermissionsRequest) (*GetAllUsersAPIPermissionsResponse, error)
	GetAllUsersViewPermissions(request *GetAllUsersViewPermissionsRequest) (*GetAllUsersViewPermissionsResponse, error)
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

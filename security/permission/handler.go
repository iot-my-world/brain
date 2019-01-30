package permission

import (
	"gitlab.com/iotTracker/brain/search"
	"gitlab.com/iotTracker/brain/security"
)

type Handler interface {
	UserHasPermission(request *UserHasPermissionRequest, response *UserHasPermissionResponse) error
	GetAllUsersPermissions(request *GetAllUsersPermissionsRequest, response *GetAllUsersPermissionsResponse) error
}

type UserHasPermissionRequest struct {
	UserIdentifier search.Identifier   `json:"userIdentifier"`
	Permission     security.Permission `json:"permission"`
}

type UserHasPermissionResponse struct {
	Result bool `json:"result"`
}

type GetAllUsersPermissionsRequest struct {
	UserIdentifier search.Identifier `json:"userIdentifier"`
}

type GetAllUsersPermissionsResponse struct {
	Permissions []security.Permission `json:"permissions"`
}

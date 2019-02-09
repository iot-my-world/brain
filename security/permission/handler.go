package permission

import (
	"gitlab.com/iotTracker/brain/search"
)

type Handler interface {
	UserHasPermission(request *UserHasPermissionRequest, response *UserHasPermissionResponse) error
	GetAllUsersPermissions(request *GetAllUsersPermissionsRequest, response *GetAllUsersPermissionsResponse) error
}

type UserHasPermissionRequest struct {
	UserIdentifier search.Identifier `json:"userIdentifier"`
	Permission     Permission        `json:"permission"`
}

type UserHasPermissionResponse struct {
	Result bool `json:"result"`
}

type GetAllUsersPermissionsRequest struct {
	UserIdentifier search.Identifier `json:"userIdentifier"`
}

type GetAllUsersPermissionsResponse struct {
	Permissions []Permission `json:"permissions"`
}

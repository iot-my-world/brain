package handler

import (
	"gitlab.com/iotTracker/brain/search"
	"gitlab.com/iotTracker/brain/security/permission"
)

type Handler interface {
	UserHasPermission(request *UserHasPermissionRequest, response *UserHasPermissionResponse) error
	GetAllUsersPermissions(request *GetAllUsersPermissionsRequest, response *GetAllUsersPermissionsResponse) error
}

type UserHasPermissionRequest struct {
	UserIdentifier search.Identifier     `json:"userIdentifier"`
	Permission     permission.Permission `json:"permission"`
}

type UserHasPermissionResponse struct {
	Result bool `json:"result"`
}

type GetAllUsersPermissionsRequest struct {
	UserIdentifier search.Identifier `json:"userIdentifier"`
}

type GetAllUsersPermissionsResponse struct {
	Permissions []permission.Permission `json:"permissions"`
}

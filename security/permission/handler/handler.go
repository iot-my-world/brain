package handler

import (
	"gitlab.com/iotTracker/brain/search/identifier"
	"gitlab.com/iotTracker/brain/security/permission/api"
	"gitlab.com/iotTracker/brain/security/permission/view"
)

type Handler interface {
	UserHasPermission(request *UserHasPermissionRequest, response *UserHasPermissionResponse) error
	GetAllUsersAPIPermissions(request *GetAllUsersAPIPermissionsRequest, response *GetAllUsersAPIPermissionsResponse) error
	GetAllUsersViewPermissions(request *GetAllUsersViewPermissionsRequest, response *GetAllUsersViewPermissionsResponse) error
}

type UserHasPermissionRequest struct {
	UserIdentifier identifier.Identifier `json:"userIdentifier"`
	Permission     api.Permission        `json:"permission"`
}

type UserHasPermissionResponse struct {
	Result bool `json:"result"`
}

type GetAllUsersAPIPermissionsRequest struct {
	UserIdentifier identifier.Identifier `json:"userIdentifier"`
}

type GetAllUsersAPIPermissionsResponse struct {
	Permissions []api.Permission `json:"permissions"`
}

type GetAllUsersViewPermissionsRequest struct {
	UserIdentifier identifier.Identifier `json:"userIdentifier"`
}

type GetAllUsersViewPermissionsResponse struct {
	Permissions []view.Permission `json:"permissions"`
}

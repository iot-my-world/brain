package handler

import (
	"gitlab.com/iotTracker/brain/search/identifier"
	"gitlab.com/iotTracker/brain/security/permission/api"
	"gitlab.com/iotTracker/brain/security/permission/view"
	"gitlab.com/iotTracker/brain/security/claims"
)

type Handler interface {
	UserHasPermission(request *UserHasPermissionRequest, response *UserHasPermissionResponse) error
	GetAllUsersAPIPermissions(request *GetAllUsersAPIPermissionsRequest, response *GetAllUsersAPIPermissionsResponse) error
	GetAllUsersViewPermissions(request *GetAllUsersViewPermissionsRequest, response *GetAllUsersViewPermissionsResponse) error
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

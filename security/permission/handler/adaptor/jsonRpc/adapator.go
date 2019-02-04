package jsonRpc

import (
	"net/http"
	"gitlab.com/iotTracker/brain/search"
	"gitlab.com/iotTracker/brain/security/permission"
	"gitlab.com/iotTracker/brain/security"
)

type service struct {
	permissionHandler permission.Handler
}

func NewService(permissionHandler permission.Handler) *service {
	return &service{
		permissionHandler: permissionHandler,
	}
}

type GetAllUsersPermissionsRequest struct {
	UserIdentifier search.IdentifierWrapper `json:"userIdentifier"`
}

type GetAllUsersPermissionsResponse struct {
	Permissions []security.Permission `json:"permission"`
}

func (s *service) GetAllUsersPermissions(r *http.Request, request *GetAllUsersPermissionsRequest, response *GetAllUsersPermissionsResponse) error {
	return nil
}

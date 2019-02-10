package jsonRpc

import (
	"gitlab.com/iotTracker/brain/search/wrappedIdentifier"
	"gitlab.com/iotTracker/brain/security/permission"
	permissionHandler "gitlab.com/iotTracker/brain/security/permission/handler"
	"net/http"
)

type adaptor struct {
	permissionHandler permissionHandler.Handler
}

func New(permissionHandler permissionHandler.Handler) *adaptor {
	return &adaptor{
		permissionHandler: permissionHandler,
	}
}

type GetAllUsersPermissionsRequest struct {
	UserIdentifier wrappedIdentifier.WrappedIdentifier `json:"userIdentifier"`
}

type GetAllUsersPermissionsResponse struct {
	Permissions []permission.Permission `json:"permission"`
}

func (s *adaptor) GetAllUsersPermissions(r *http.Request, request *GetAllUsersPermissionsRequest, response *GetAllUsersPermissionsResponse) error {
	id, err := request.UserIdentifier.UnWrap()
	if err != nil {
		return err
	}

	getAllUsersPermissionsResponse := permissionHandler.GetAllUsersPermissionsResponse{}
	if err := s.permissionHandler.GetAllUsersPermissions(&permissionHandler.GetAllUsersPermissionsRequest{
		UserIdentifier: id,
	}, &getAllUsersPermissionsResponse); err != nil {
		return err
	}
	response.Permissions = getAllUsersPermissionsResponse.Permissions
	return nil
}

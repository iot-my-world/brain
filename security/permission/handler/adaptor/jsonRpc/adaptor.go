package jsonRpc

import (
	"gitlab.com/iotTracker/brain/search/wrappedIdentifier"
	permissionHandler "gitlab.com/iotTracker/brain/security/permission/handler"
	"net/http"
	"gitlab.com/iotTracker/brain/security/permission/api"
	"gitlab.com/iotTracker/brain/security/permission/view"
)

type adaptor struct {
	permissionHandler permissionHandler.Handler
}

func New(permissionHandler permissionHandler.Handler) *adaptor {
	return &adaptor{
		permissionHandler: permissionHandler,
	}
}

type GetAllUsersAPIPermissionsRequest struct {
	UserIdentifier wrappedIdentifier.WrappedIdentifier `json:"userIdentifier"`
}

type GetAllUsersAPIPermissionsResponse struct {
	Permissions []api.Permission `json:"permission"`
}

func (s *adaptor) GetAllUsersAPIPermissions(r *http.Request, request *GetAllUsersAPIPermissionsRequest, response *GetAllUsersAPIPermissionsResponse) error {
	id, err := request.UserIdentifier.UnWrap()
	if err != nil {
		return err
	}

	getAllUsersAPIPermissionsResponse := permissionHandler.GetAllUsersAPIPermissionsResponse{}
	if err := s.permissionHandler.GetAllUsersAPIPermissions(&permissionHandler.GetAllUsersAPIPermissionsRequest{
		UserIdentifier: id,
	}, &getAllUsersAPIPermissionsResponse); err != nil {
		return err
	}
	response.Permissions = getAllUsersAPIPermissionsResponse.Permissions
	return nil
}

type GetAllUsersViewPermissionsRequest struct {
	UserIdentifier wrappedIdentifier.WrappedIdentifier `json:"userIdentifier"`
}

type GetAllUsersViewPermissionsResponse struct {
	Permissions []view.Permission `json:"permission"`
}

func (s *adaptor) GetAllUsersViewPermissions(r *http.Request, request *GetAllUsersViewPermissionsRequest, response *GetAllUsersViewPermissionsResponse) error {
	id, err := request.UserIdentifier.UnWrap()
	if err != nil {
		return err
	}

	getAllUsersViewPermissionsResponse := permissionHandler.GetAllUsersViewPermissionsResponse{}
	if err := s.permissionHandler.GetAllUsersViewPermissions(&permissionHandler.GetAllUsersViewPermissionsRequest{
		UserIdentifier: id,
	}, &getAllUsersViewPermissionsResponse); err != nil {
		return err
	}
	response.Permissions = getAllUsersViewPermissionsResponse.Permissions
	return nil
}

package jsonRpc

import (
	"gitlab.com/iotTracker/brain/search"
	"gitlab.com/iotTracker/brain/security/permission"
	"net/http"
)

type adaptor struct {
	permissionHandler permission.Handler
}

func New(permissionHandler permission.Handler) *adaptor {
	return &adaptor{
		permissionHandler: permissionHandler,
	}
}

type GetAllUsersPermissionsRequest struct {
	UserIdentifier search.WrappedIdentifier `json:"userIdentifier"`
}

type GetAllUsersPermissionsResponse struct {
	Permissions []permission.Permission `json:"permission"`
}

func (s *adaptor) GetAllUsersPermissions(r *http.Request, request *GetAllUsersPermissionsRequest, response *GetAllUsersPermissionsResponse) error {
	id, err := request.UserIdentifier.UnWrap()
	if err != nil {
		return err
	}

	getAllUsersPermissionsResponse := permission.GetAllUsersPermissionsResponse{}
	if err := s.permissionHandler.GetAllUsersPermissions(&permission.GetAllUsersPermissionsRequest{
		UserIdentifier: id,
	}, &getAllUsersPermissionsResponse); err != nil {
		return err
	}
	response.Permissions = getAllUsersPermissionsResponse.Permissions
	return nil
}

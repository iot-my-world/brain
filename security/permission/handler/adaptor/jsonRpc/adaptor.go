package jsonRpc

import (
	"gitlab.com/iotTracker/brain/search/wrappedIdentifier"
	"gitlab.com/iotTracker/brain/security/permission/api"
	permissionHandler "gitlab.com/iotTracker/brain/security/permission/handler"
	"gitlab.com/iotTracker/brain/security/permission/view"
	"net/http"
	"gitlab.com/iotTracker/brain/security/wrappedClaims"
	"gitlab.com/iotTracker/brain/log"
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
	claims, err := wrappedClaims.UnwrapClaimsFromContext(r)
	if err != nil {
		log.Warn(err.Error())
		return err
	}

	id, err := request.UserIdentifier.UnWrap()
	if err != nil {
		return err
	}

	getAllUsersAPIPermissionsResponse := permissionHandler.GetAllUsersAPIPermissionsResponse{}
	if err := s.permissionHandler.GetAllUsersAPIPermissions(&permissionHandler.GetAllUsersAPIPermissionsRequest{
		Claims:         claims,
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
	claims, err := wrappedClaims.UnwrapClaimsFromContext(r)
	if err != nil {
		log.Warn(err.Error())
		return err
	}

	id, err := request.UserIdentifier.UnWrap()
	if err != nil {
		return err
	}

	getAllUsersViewPermissionsResponse := permissionHandler.GetAllUsersViewPermissionsResponse{}
	if err := s.permissionHandler.GetAllUsersViewPermissions(&permissionHandler.GetAllUsersViewPermissionsRequest{
		Claims:         claims,
		UserIdentifier: id,
	}, &getAllUsersViewPermissionsResponse); err != nil {
		return err
	}
	response.Permissions = getAllUsersViewPermissionsResponse.Permissions
	return nil
}

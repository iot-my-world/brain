package jsonRpc

import (
	"gitlab.com/iotTracker/brain/log"
	"gitlab.com/iotTracker/brain/search/wrappedIdentifier"
	permissionAdministrator "gitlab.com/iotTracker/brain/security/permission/administrator"
	"gitlab.com/iotTracker/brain/security/permission/api"
	"gitlab.com/iotTracker/brain/security/permission/view"
	"gitlab.com/iotTracker/brain/security/wrappedClaims"
	"net/http"
)

type adaptor struct {
	permissionAdministrator permissionAdministrator.Handler
}

func New(permissionAdministrator permissionAdministrator.Handler) *adaptor {
	return &adaptor{
		permissionAdministrator: permissionAdministrator,
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

	getAllUsersAPIPermissionsResponse := permissionAdministrator.GetAllUsersAPIPermissionsResponse{}
	if err := s.permissionAdministrator.GetAllUsersAPIPermissions(&permissionAdministrator.GetAllUsersAPIPermissionsRequest{
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

	getAllUsersViewPermissionsResponse := permissionAdministrator.GetAllUsersViewPermissionsResponse{}
	if err := s.permissionAdministrator.GetAllUsersViewPermissions(&permissionAdministrator.GetAllUsersViewPermissionsRequest{
		Claims:         claims,
		UserIdentifier: id,
	}, &getAllUsersViewPermissionsResponse); err != nil {
		return err
	}
	response.Permissions = getAllUsersViewPermissionsResponse.Permissions
	return nil
}

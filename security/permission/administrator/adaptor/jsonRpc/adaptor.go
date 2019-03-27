package jsonRpc

import (
	"gitlab.com/iotTracker/brain/log"
	wrappedIdentifier "gitlab.com/iotTracker/brain/search/identifier/wrapped"
	wrappedClaims "gitlab.com/iotTracker/brain/security/claims/wrapped"
	permissionAdministrator "gitlab.com/iotTracker/brain/security/permission/administrator"
	"gitlab.com/iotTracker/brain/security/permission/api"
	"gitlab.com/iotTracker/brain/security/permission/view"
	"net/http"
)

type adaptor struct {
	permissionAdministrator permissionAdministrator.Administrator
}

func New(permissionAdministrator permissionAdministrator.Administrator) *adaptor {
	return &adaptor{
		permissionAdministrator: permissionAdministrator,
	}
}

type GetAllUsersAPIPermissionsRequest struct {
	UserIdentifier wrappedIdentifier.Wrapped `json:"userIdentifier"`
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

	getAllUsersAPIPermissionsResponse, err := s.permissionAdministrator.GetAllUsersAPIPermissions(&permissionAdministrator.GetAllUsersAPIPermissionsRequest{
		Claims:         claims,
		UserIdentifier: id,
	})
	if err != nil {
		return err
	}
	response.Permissions = getAllUsersAPIPermissionsResponse.Permissions
	return nil
}

type GetAllUsersViewPermissionsRequest struct {
	UserIdentifier wrappedIdentifier.Wrapped `json:"userIdentifier"`
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

	getAllUsersViewPermissionsResponse, err := s.permissionAdministrator.GetAllUsersViewPermissions(&permissionAdministrator.GetAllUsersViewPermissionsRequest{
		Claims:         claims,
		UserIdentifier: id,
	})
	if err != nil {
		return err
	}
	response.Permissions = getAllUsersViewPermissionsResponse.Permissions
	return nil
}

package jsonRpc

import (
	"github.com/iot-my-world/brain/internal/log"
	wrappedIdentifier "github.com/iot-my-world/brain/pkg/search/identifier/wrapped"
	wrappedClaims "github.com/iot-my-world/brain/security/claims/wrapped"
	permissionAdministrator "github.com/iot-my-world/brain/security/permission/administrator"
	"github.com/iot-my-world/brain/security/permission/api"
	"github.com/iot-my-world/brain/security/permission/view"
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
	WrappedUserIdentifier wrappedIdentifier.Wrapped `json:"userIdentifier"`
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

	getAllUsersAPIPermissionsResponse, err := s.permissionAdministrator.GetAllUsersAPIPermissions(&permissionAdministrator.GetAllUsersAPIPermissionsRequest{
		Claims:         claims,
		UserIdentifier: request.WrappedUserIdentifier.Identifier,
	})
	if err != nil {
		return err
	}
	response.Permissions = getAllUsersAPIPermissionsResponse.Permissions
	return nil
}

type GetAllUsersViewPermissionsRequest struct {
	WrappedUserIdentifier wrappedIdentifier.Wrapped `json:"userIdentifier"`
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

	getAllUsersViewPermissionsResponse, err := s.permissionAdministrator.GetAllUsersViewPermissions(&permissionAdministrator.GetAllUsersViewPermissionsRequest{
		Claims:         claims,
		UserIdentifier: request.WrappedUserIdentifier.Identifier,
	})
	if err != nil {
		return err
	}
	response.Permissions = getAllUsersViewPermissionsResponse.Permissions
	return nil
}

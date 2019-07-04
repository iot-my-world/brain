package jsonRpc

import (
	"github.com/iot-my-world/brain/internal/log"
	wrappedIdentifier "github.com/iot-my-world/brain/pkg/search/identifier/wrapped"
	"github.com/iot-my-world/brain/pkg/security/claims/wrapped"
	"github.com/iot-my-world/brain/pkg/security/permission/administrator"
	api2 "github.com/iot-my-world/brain/pkg/security/permission/api"
	view2 "github.com/iot-my-world/brain/pkg/security/permission/view"
	"net/http"
)

type adaptor struct {
	permissionAdministrator administrator.Administrator
}

func New(permissionAdministrator administrator.Administrator) *adaptor {
	return &adaptor{
		permissionAdministrator: permissionAdministrator,
	}
}

type GetAllUsersAPIPermissionsRequest struct {
	WrappedUserIdentifier wrappedIdentifier.Wrapped `json:"userIdentifier"`
}

type GetAllUsersAPIPermissionsResponse struct {
	Permissions []api2.Permission `json:"permission"`
}

func (s *adaptor) GetAllUsersAPIPermissions(r *http.Request, request *GetAllUsersAPIPermissionsRequest, response *GetAllUsersAPIPermissionsResponse) error {
	claims, err := wrapped.UnwrapClaimsFromContext(r)
	if err != nil {
		log.Warn(err.Error())
		return err
	}

	getAllUsersAPIPermissionsResponse, err := s.permissionAdministrator.GetAllUsersAPIPermissions(&administrator.GetAllUsersAPIPermissionsRequest{
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
	Permissions []view2.Permission `json:"permission"`
}

func (s *adaptor) GetAllUsersViewPermissions(r *http.Request, request *GetAllUsersViewPermissionsRequest, response *GetAllUsersViewPermissionsResponse) error {
	claims, err := wrapped.UnwrapClaimsFromContext(r)
	if err != nil {
		log.Warn(err.Error())
		return err
	}

	getAllUsersViewPermissionsResponse, err := s.permissionAdministrator.GetAllUsersViewPermissions(&administrator.GetAllUsersViewPermissionsRequest{
		Claims:         claims,
		UserIdentifier: request.WrappedUserIdentifier.Identifier,
	})
	if err != nil {
		return err
	}
	response.Permissions = getAllUsersViewPermissionsResponse.Permissions
	return nil
}

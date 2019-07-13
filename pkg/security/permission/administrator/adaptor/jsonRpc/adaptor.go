package jsonRpc

import (
	"github.com/iot-my-world/brain/internal/log"
	jsonRpcServiceProvider "github.com/iot-my-world/brain/pkg/api/jsonRpc/service/provider"
	wrappedIdentifier "github.com/iot-my-world/brain/pkg/search/identifier/wrapped"
	"github.com/iot-my-world/brain/pkg/security/claims/wrapped"
	"github.com/iot-my-world/brain/pkg/security/permission/administrator"
	apiPermission "github.com/iot-my-world/brain/pkg/security/permission/api"
	viewPermission "github.com/iot-my-world/brain/pkg/security/permission/view"
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

func (a *adaptor) Name() jsonRpcServiceProvider.Name {
	return jsonRpcServiceProvider.Name(administrator.ServiceProvider)
}

func (a *adaptor) MethodRequiresAuthorization(string) bool {
	return true
}

type GetAllUsersAPIPermissionsRequest struct {
	WrappedUserIdentifier wrappedIdentifier.Wrapped `json:"userIdentifier"`
}

type GetAllUsersAPIPermissionsResponse struct {
	Permissions []apiPermission.Permission `json:"permission"`
}

func (a *adaptor) GetAllUsersAPIPermissions(r *http.Request, request *GetAllUsersAPIPermissionsRequest, response *GetAllUsersAPIPermissionsResponse) error {
	claims, err := wrapped.UnwrapClaimsFromContext(r)
	if err != nil {
		log.Warn(err.Error())
		return err
	}

	getAllUsersAPIPermissionsResponse, err := a.permissionAdministrator.GetAllUsersAPIPermissions(&administrator.GetAllUsersAPIPermissionsRequest{
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
	Permissions []viewPermission.Permission `json:"permission"`
}

func (a *adaptor) GetAllUsersViewPermissions(r *http.Request, request *GetAllUsersViewPermissionsRequest, response *GetAllUsersViewPermissionsResponse) error {
	claims, err := wrapped.UnwrapClaimsFromContext(r)
	if err != nil {
		log.Warn(err.Error())
		return err
	}

	getAllUsersViewPermissionsResponse, err := a.permissionAdministrator.GetAllUsersViewPermissions(&administrator.GetAllUsersViewPermissionsRequest{
		Claims:         claims,
		UserIdentifier: request.WrappedUserIdentifier.Identifier,
	})
	if err != nil {
		return err
	}
	response.Permissions = getAllUsersViewPermissionsResponse.Permissions
	return nil
}

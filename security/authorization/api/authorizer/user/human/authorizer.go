package human

import (
	brainException "github.com/iot-my-world/brain/internal/exception"
	apiAuthorizer "github.com/iot-my-world/brain/security/authorization/api/authorizer"
	apiAuthException "github.com/iot-my-world/brain/security/authorization/api/authorizer/exception"
	humanUserLoginClaims "github.com/iot-my-world/brain/security/claims/login/user/human"
	"github.com/iot-my-world/brain/security/claims/registerClientAdminUser"
	"github.com/iot-my-world/brain/security/claims/registerClientUser"
	"github.com/iot-my-world/brain/security/claims/registerCompanyAdminUser"
	"github.com/iot-my-world/brain/security/claims/registerCompanyUser"
	"github.com/iot-my-world/brain/security/claims/resetPassword"
	wrappedClaims "github.com/iot-my-world/brain/security/claims/wrapped"
	permissionAdministrator "github.com/iot-my-world/brain/security/permission/administrator"
	"github.com/iot-my-world/brain/security/permission/api"
	"github.com/iot-my-world/brain/security/token"
)

type authorizer struct {
	jwtValidator      token.JWTValidator
	permissionHandler permissionAdministrator.Administrator
}

func New(
	jwtValidator token.JWTValidator,
	permissionHandler permissionAdministrator.Administrator,
) apiAuthorizer.Authorizer {
	return &authorizer{
		jwtValidator:      jwtValidator,
		permissionHandler: permissionHandler,
	}
}

func (a *authorizer) AuthorizeAPIReq(jwt string, jsonRpcMethod string) (wrappedClaims.Wrapped, error) {

	// Validate the jwt
	wrappedJWTClaims, err := a.jwtValidator.ValidateJWT(jwt)
	if err != nil {
		return wrappedClaims.Wrapped{}, err
	}
	unwrappedJWTClaims, err := wrappedJWTClaims.Unwrap()
	if err != nil {
		return wrappedClaims.Wrapped{}, err
	}

	switch typedClaims := unwrappedJWTClaims.(type) {
	case humanUserLoginClaims.Login:
		// if these are login claims we check in the normal way if the user has the
		// required permission to check access the api
		userHasPermissionResponse, err := a.permissionHandler.UserHasPermission(&permissionAdministrator.UserHasPermissionRequest{
			Claims:         typedClaims,
			UserIdentifier: typedClaims.UserId,
			Permission:     api.Permission(jsonRpcMethod),
		})
		if err != nil {
			return wrappedClaims.Wrapped{}, brainException.Unexpected{Reasons: []string{"determining if user has permission", err.Error()}}
		}
		if !userHasPermissionResponse.Result {
			return wrappedClaims.Wrapped{}, apiAuthException.NotAuthorised{Permission: api.Permission(jsonRpcMethod)}
		}
		// user was authorised
		return wrappedJWTClaims, nil

	case registerCompanyAdminUser.RegisterCompanyAdminUser:
		permissionForMethod := api.Permission(jsonRpcMethod)
		// check the permissions granted by the RegisterCompanyAdminUser claims to see if this
		// method is allowed
		for allowedPermIdx := range registerCompanyAdminUser.GrantedAPIPermissions {
			if registerCompanyAdminUser.GrantedAPIPermissions[allowedPermIdx] == permissionForMethod {
				return wrappedJWTClaims, nil
			}
			if allowedPermIdx == len(registerCompanyAdminUser.GrantedAPIPermissions)-1 {
				return wrappedClaims.Wrapped{}, apiAuthException.NotAuthorised{Permission: api.Permission(jsonRpcMethod)}
			}
		}

	case registerCompanyUser.RegisterCompanyUser:
		permissionForMethod := api.Permission(jsonRpcMethod)
		// check the permissions granted by the RegisterCompanyUser claims to see if this
		// method is allowed
		for allowedPermIdx := range registerCompanyUser.GrantedAPIPermissions {
			if registerCompanyUser.GrantedAPIPermissions[allowedPermIdx] == permissionForMethod {
				return wrappedJWTClaims, nil
			}
			if allowedPermIdx == len(registerCompanyUser.GrantedAPIPermissions)-1 {
				return wrappedClaims.Wrapped{}, apiAuthException.NotAuthorised{Permission: api.Permission(jsonRpcMethod)}
			}
		}

	case registerClientAdminUser.RegisterClientAdminUser:
		permissionForMethod := api.Permission(jsonRpcMethod)
		// check the permissions granted by the RegisterClientAdminUser claims to see if this
		// method is allowed
		for allowedPermIdx := range registerClientAdminUser.GrantedAPIPermissions {
			if registerClientAdminUser.GrantedAPIPermissions[allowedPermIdx] == permissionForMethod {
				return wrappedJWTClaims, nil
			}
			if allowedPermIdx == len(registerClientAdminUser.GrantedAPIPermissions)-1 {
				return wrappedClaims.Wrapped{}, apiAuthException.NotAuthorised{Permission: api.Permission(jsonRpcMethod)}
			}
		}

	case registerClientUser.RegisterClientUser:
		permissionForMethod := api.Permission(jsonRpcMethod)
		// check the permissions granted by the RegisterClientUser claims to see if this
		// method is allowed
		for allowedPermIdx := range registerClientUser.GrantedAPIPermissions {
			if registerClientUser.GrantedAPIPermissions[allowedPermIdx] == permissionForMethod {
				return wrappedJWTClaims, nil
			}
			if allowedPermIdx == len(registerClientUser.GrantedAPIPermissions)-1 {
				return wrappedClaims.Wrapped{}, apiAuthException.NotAuthorised{Permission: api.Permission(jsonRpcMethod)}
			}
		}

	case resetPassword.ResetPassword:
		permissionForMethod := api.Permission(jsonRpcMethod)
		// check the permissions granted by the ResetPassword claims to see if this
		// method is allowed
		for allowedPermIdx := range resetPassword.GrantedAPIPermissions {
			if resetPassword.GrantedAPIPermissions[allowedPermIdx] == permissionForMethod {
				return wrappedJWTClaims, nil
			}
			if allowedPermIdx == len(resetPassword.GrantedAPIPermissions)-1 {
				return wrappedClaims.Wrapped{}, apiAuthException.NotAuthorised{Permission: api.Permission(jsonRpcMethod)}
			}
		}

	default:
		return wrappedClaims.Wrapped{}, apiAuthException.NotAuthorised{Permission: api.Permission(jsonRpcMethod)}
	}

	return wrappedClaims.Wrapped{}, apiAuthException.NotAuthorised{Permission: api.Permission(jsonRpcMethod)}
}

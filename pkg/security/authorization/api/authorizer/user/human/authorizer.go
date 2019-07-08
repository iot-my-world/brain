package human

import (
	brainException "github.com/iot-my-world/brain/internal/exception"
	apiAuthorizer "github.com/iot-my-world/brain/pkg/security/authorization/api/authorizer"
	"github.com/iot-my-world/brain/pkg/security/authorization/api/authorizer/exception"
	"github.com/iot-my-world/brain/pkg/security/claims/login/user/human"
	registerClientAdminUserClaims "github.com/iot-my-world/brain/pkg/security/claims/registerClientAdminUser"
	registerClientUserClaims "github.com/iot-my-world/brain/pkg/security/claims/registerClientUser"
	registerCompanyAdminUserClaims "github.com/iot-my-world/brain/pkg/security/claims/registerCompanyAdminUser"
	registerCompanyUserClaims "github.com/iot-my-world/brain/pkg/security/claims/registerCompanyUser"
	resetPasswordClaims "github.com/iot-my-world/brain/pkg/security/claims/resetPassword"
	"github.com/iot-my-world/brain/pkg/security/claims/wrapped"
	"github.com/iot-my-world/brain/pkg/security/permission/administrator"
	apiPermissions "github.com/iot-my-world/brain/pkg/security/permission/api"
	"github.com/iot-my-world/brain/pkg/security/token"
)

type authorizer struct {
	jwtValidator      token.JWTValidator
	permissionHandler administrator.Administrator
}

func New(
	jwtValidator token.JWTValidator,
	permissionHandler administrator.Administrator,
) apiAuthorizer.Authorizer {
	return &authorizer{
		jwtValidator:      jwtValidator,
		permissionHandler: permissionHandler,
	}
}

func (a *authorizer) AuthorizeAPIReq(jwt string, jsonRpcMethod string) (wrapped.Wrapped, error) {

	// Validate the jwt
	wrappedJWTClaims, err := a.jwtValidator.ValidateJWT(jwt)
	if err != nil {
		return wrapped.Wrapped{}, err
	}
	unwrappedJWTClaims, err := wrappedJWTClaims.Unwrap()
	if err != nil {
		return wrapped.Wrapped{}, err
	}

	switch typedClaims := unwrappedJWTClaims.(type) {
	case human.Login:
		// if these are login claims we check in the normal way if the user has the
		// required permission to check access the api
		userHasPermissionResponse, err := a.permissionHandler.UserHasPermission(&administrator.UserHasPermissionRequest{
			Claims:         typedClaims,
			UserIdentifier: typedClaims.UserId,
			Permission:     apiPermissions.Permission(jsonRpcMethod),
		})
		if err != nil {
			return wrapped.Wrapped{}, brainException.Unexpected{Reasons: []string{"determining if user has permission", err.Error()}}
		}
		if !userHasPermissionResponse.Result {
			return wrapped.Wrapped{}, exception.NotAuthorised{Permission: apiPermissions.Permission(jsonRpcMethod)}
		}
		// user was authorised
		return wrappedJWTClaims, nil

	case registerCompanyAdminUserClaims.RegisterCompanyAdminUser:
		permissionForMethod := apiPermissions.Permission(jsonRpcMethod)
		// check the permissions granted by the RegisterCompanyAdminUser claims to see if this
		// method is allowed
		for allowedPermIdx := range registerCompanyAdminUserClaims.GrantedAPIPermissions {
			if registerCompanyAdminUserClaims.GrantedAPIPermissions[allowedPermIdx] == permissionForMethod {
				return wrappedJWTClaims, nil
			}
			if allowedPermIdx == len(registerCompanyAdminUserClaims.GrantedAPIPermissions)-1 {
				return wrapped.Wrapped{}, exception.NotAuthorised{Permission: apiPermissions.Permission(jsonRpcMethod)}
			}
		}

	case registerCompanyUserClaims.RegisterCompanyUser:
		permissionForMethod := apiPermissions.Permission(jsonRpcMethod)
		// check the permissions granted by the RegisterCompanyUser claims to see if this
		// method is allowed
		for allowedPermIdx := range registerCompanyUserClaims.GrantedAPIPermissions {
			if registerCompanyUserClaims.GrantedAPIPermissions[allowedPermIdx] == permissionForMethod {
				return wrappedJWTClaims, nil
			}
			if allowedPermIdx == len(registerCompanyUserClaims.GrantedAPIPermissions)-1 {
				return wrapped.Wrapped{}, exception.NotAuthorised{Permission: apiPermissions.Permission(jsonRpcMethod)}
			}
		}

	case registerClientAdminUserClaims.RegisterClientAdminUser:
		permissionForMethod := apiPermissions.Permission(jsonRpcMethod)
		// check the permissions granted by the RegisterClientAdminUser claims to see if this
		// method is allowed
		for allowedPermIdx := range registerClientAdminUserClaims.GrantedAPIPermissions {
			if registerClientAdminUserClaims.GrantedAPIPermissions[allowedPermIdx] == permissionForMethod {
				return wrappedJWTClaims, nil
			}
			if allowedPermIdx == len(registerClientAdminUserClaims.GrantedAPIPermissions)-1 {
				return wrapped.Wrapped{}, exception.NotAuthorised{Permission: apiPermissions.Permission(jsonRpcMethod)}
			}
		}

	case registerClientUserClaims.RegisterClientUser:
		permissionForMethod := apiPermissions.Permission(jsonRpcMethod)
		// check the permissions granted by the RegisterClientUser claims to see if this
		// method is allowed
		for allowedPermIdx := range registerClientUserClaims.GrantedAPIPermissions {
			if registerClientUserClaims.GrantedAPIPermissions[allowedPermIdx] == permissionForMethod {
				return wrappedJWTClaims, nil
			}
			if allowedPermIdx == len(registerClientUserClaims.GrantedAPIPermissions)-1 {
				return wrapped.Wrapped{}, exception.NotAuthorised{Permission: apiPermissions.Permission(jsonRpcMethod)}
			}
		}

	case resetPasswordClaims.ResetPassword:
		permissionForMethod := apiPermissions.Permission(jsonRpcMethod)
		// check the permissions granted by the ResetPassword claims to see if this
		// method is allowed
		for allowedPermIdx := range resetPasswordClaims.GrantedAPIPermissions {
			if resetPasswordClaims.GrantedAPIPermissions[allowedPermIdx] == permissionForMethod {
				return wrappedJWTClaims, nil
			}
			if allowedPermIdx == len(resetPasswordClaims.GrantedAPIPermissions)-1 {
				return wrapped.Wrapped{}, exception.NotAuthorised{Permission: apiPermissions.Permission(jsonRpcMethod)}
			}
		}

	default:
		return wrapped.Wrapped{}, exception.NotAuthorised{Permission: apiPermissions.Permission(jsonRpcMethod)}
	}

	return wrapped.Wrapped{}, exception.NotAuthorised{Permission: apiPermissions.Permission(jsonRpcMethod)}
}

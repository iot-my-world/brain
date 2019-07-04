package human

import (
	brainException "github.com/iot-my-world/brain/internal/exception"
	authorizer2 "github.com/iot-my-world/brain/pkg/security/authorization/api/authorizer"
	"github.com/iot-my-world/brain/pkg/security/authorization/api/authorizer/exception"
	"github.com/iot-my-world/brain/pkg/security/claims/login/user/human"
	registerClientAdminUser2 "github.com/iot-my-world/brain/pkg/security/claims/registerClientAdminUser"
	registerClientUser2 "github.com/iot-my-world/brain/pkg/security/claims/registerClientUser"
	registerCompanyAdminUser2 "github.com/iot-my-world/brain/pkg/security/claims/registerCompanyAdminUser"
	registerCompanyUser2 "github.com/iot-my-world/brain/pkg/security/claims/registerCompanyUser"
	resetPassword2 "github.com/iot-my-world/brain/pkg/security/claims/resetPassword"
	"github.com/iot-my-world/brain/pkg/security/claims/wrapped"
	"github.com/iot-my-world/brain/pkg/security/permission/administrator"
	api2 "github.com/iot-my-world/brain/pkg/security/permission/api"
	token2 "github.com/iot-my-world/brain/pkg/security/token"
)

type authorizer struct {
	jwtValidator      token2.JWTValidator
	permissionHandler administrator.Administrator
}

func New(
	jwtValidator token2.JWTValidator,
	permissionHandler administrator.Administrator,
) authorizer2.Authorizer {
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
			Permission:     api2.Permission(jsonRpcMethod),
		})
		if err != nil {
			return wrapped.Wrapped{}, brainException.Unexpected{Reasons: []string{"determining if user has permission", err.Error()}}
		}
		if !userHasPermissionResponse.Result {
			return wrapped.Wrapped{}, exception.NotAuthorised{Permission: api2.Permission(jsonRpcMethod)}
		}
		// user was authorised
		return wrappedJWTClaims, nil

	case registerCompanyAdminUser2.RegisterCompanyAdminUser:
		permissionForMethod := api2.Permission(jsonRpcMethod)
		// check the permissions granted by the RegisterCompanyAdminUser claims to see if this
		// method is allowed
		for allowedPermIdx := range registerCompanyAdminUser2.GrantedAPIPermissions {
			if registerCompanyAdminUser2.GrantedAPIPermissions[allowedPermIdx] == permissionForMethod {
				return wrappedJWTClaims, nil
			}
			if allowedPermIdx == len(registerCompanyAdminUser2.GrantedAPIPermissions)-1 {
				return wrapped.Wrapped{}, exception.NotAuthorised{Permission: api2.Permission(jsonRpcMethod)}
			}
		}

	case registerCompanyUser2.RegisterCompanyUser:
		permissionForMethod := api2.Permission(jsonRpcMethod)
		// check the permissions granted by the RegisterCompanyUser claims to see if this
		// method is allowed
		for allowedPermIdx := range registerCompanyUser2.GrantedAPIPermissions {
			if registerCompanyUser2.GrantedAPIPermissions[allowedPermIdx] == permissionForMethod {
				return wrappedJWTClaims, nil
			}
			if allowedPermIdx == len(registerCompanyUser2.GrantedAPIPermissions)-1 {
				return wrapped.Wrapped{}, exception.NotAuthorised{Permission: api2.Permission(jsonRpcMethod)}
			}
		}

	case registerClientAdminUser2.RegisterClientAdminUser:
		permissionForMethod := api2.Permission(jsonRpcMethod)
		// check the permissions granted by the RegisterClientAdminUser claims to see if this
		// method is allowed
		for allowedPermIdx := range registerClientAdminUser2.GrantedAPIPermissions {
			if registerClientAdminUser2.GrantedAPIPermissions[allowedPermIdx] == permissionForMethod {
				return wrappedJWTClaims, nil
			}
			if allowedPermIdx == len(registerClientAdminUser2.GrantedAPIPermissions)-1 {
				return wrapped.Wrapped{}, exception.NotAuthorised{Permission: api2.Permission(jsonRpcMethod)}
			}
		}

	case registerClientUser2.RegisterClientUser:
		permissionForMethod := api2.Permission(jsonRpcMethod)
		// check the permissions granted by the RegisterClientUser claims to see if this
		// method is allowed
		for allowedPermIdx := range registerClientUser2.GrantedAPIPermissions {
			if registerClientUser2.GrantedAPIPermissions[allowedPermIdx] == permissionForMethod {
				return wrappedJWTClaims, nil
			}
			if allowedPermIdx == len(registerClientUser2.GrantedAPIPermissions)-1 {
				return wrapped.Wrapped{}, exception.NotAuthorised{Permission: api2.Permission(jsonRpcMethod)}
			}
		}

	case resetPassword2.ResetPassword:
		permissionForMethod := api2.Permission(jsonRpcMethod)
		// check the permissions granted by the ResetPassword claims to see if this
		// method is allowed
		for allowedPermIdx := range resetPassword2.GrantedAPIPermissions {
			if resetPassword2.GrantedAPIPermissions[allowedPermIdx] == permissionForMethod {
				return wrappedJWTClaims, nil
			}
			if allowedPermIdx == len(resetPassword2.GrantedAPIPermissions)-1 {
				return wrapped.Wrapped{}, exception.NotAuthorised{Permission: api2.Permission(jsonRpcMethod)}
			}
		}

	default:
		return wrapped.Wrapped{}, exception.NotAuthorised{Permission: api2.Permission(jsonRpcMethod)}
	}

	return wrapped.Wrapped{}, exception.NotAuthorised{Permission: api2.Permission(jsonRpcMethod)}
}

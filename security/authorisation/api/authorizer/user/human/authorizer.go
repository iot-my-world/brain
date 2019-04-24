package human

import (
	brainException "gitlab.com/iotTracker/brain/exception"
	apiAuthorizer "gitlab.com/iotTracker/brain/security/authorisation/api/authorizer"
	apiAuthException "gitlab.com/iotTracker/brain/security/authorisation/api/authorizer/exception"
	"gitlab.com/iotTracker/brain/security/claims/login"
	"gitlab.com/iotTracker/brain/security/claims/registerClientAdminUser"
	"gitlab.com/iotTracker/brain/security/claims/registerClientUser"
	"gitlab.com/iotTracker/brain/security/claims/registerCompanyAdminUser"
	"gitlab.com/iotTracker/brain/security/claims/registerCompanyUser"
	"gitlab.com/iotTracker/brain/security/claims/resetPassword"
	wrappedClaims "gitlab.com/iotTracker/brain/security/claims/wrapped"
	permissionAdministrator "gitlab.com/iotTracker/brain/security/permission/administrator"
	"gitlab.com/iotTracker/brain/security/permission/api"
	"gitlab.com/iotTracker/brain/security/token"
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
	case login.Login:
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

package authoriser

import (
	brainException "github.com/iot-my-world/brain/internal/exception"
	jsonRpcServerAuthoriser "github.com/iot-my-world/brain/pkg/api/jsonRpc/server/authoriser"
	"github.com/iot-my-world/brain/pkg/security/authorization/api/authorizer/exception"
	humanUserLoginClaims "github.com/iot-my-world/brain/pkg/security/claims/login/user/human"
	registerClientAdminUserClaims "github.com/iot-my-world/brain/pkg/security/claims/registerClientAdminUser"
	registerClientUserClaims "github.com/iot-my-world/brain/pkg/security/claims/registerClientUser"
	registerCompanyAdminUserClaims "github.com/iot-my-world/brain/pkg/security/claims/registerCompanyAdminUser"
	registerCompanyUserClaims "github.com/iot-my-world/brain/pkg/security/claims/registerCompanyUser"
	resetPasswordClaims "github.com/iot-my-world/brain/pkg/security/claims/resetPassword"
	wrappedClaims "github.com/iot-my-world/brain/pkg/security/claims/wrapped"
	permissionAdministrator "github.com/iot-my-world/brain/pkg/security/permission/administrator"
	apiPermissions "github.com/iot-my-world/brain/pkg/security/permission/api"
	"github.com/iot-my-world/brain/pkg/security/token"
)

type authoriser struct {
	jwtValidator      token.JWTValidator
	permissionHandler permissionAdministrator.Administrator
}

func New(
	jwtValidator token.JWTValidator,
	permissionAdministrator permissionAdministrator.Administrator,
) jsonRpcServerAuthoriser.Authoriser {
	return &authoriser{}
}

func (a *authoriser) AuthoriseServiceMethod(jwt string, jsonRpcMethod string) (wrappedClaims.Wrapped, error) {

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
			Permission:     apiPermissions.Permission(jsonRpcMethod),
		})
		if err != nil {
			return wrappedClaims.Wrapped{}, brainException.Unexpected{Reasons: []string{"determining if user has permission", err.Error()}}
		}
		if !userHasPermissionResponse.Result {
			return wrappedClaims.Wrapped{}, exception.NotAuthorised{Permission: apiPermissions.Permission(jsonRpcMethod)}
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
				return wrappedClaims.Wrapped{}, exception.NotAuthorised{Permission: apiPermissions.Permission(jsonRpcMethod)}
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
				return wrappedClaims.Wrapped{}, exception.NotAuthorised{Permission: apiPermissions.Permission(jsonRpcMethod)}
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
				return wrappedClaims.Wrapped{}, exception.NotAuthorised{Permission: apiPermissions.Permission(jsonRpcMethod)}
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
				return wrappedClaims.Wrapped{}, exception.NotAuthorised{Permission: apiPermissions.Permission(jsonRpcMethod)}
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
				return wrappedClaims.Wrapped{}, exception.NotAuthorised{Permission: apiPermissions.Permission(jsonRpcMethod)}
			}
		}

	default:
		return wrappedClaims.Wrapped{}, exception.NotAuthorised{Permission: apiPermissions.Permission(jsonRpcMethod)}
	}

	return wrappedClaims.Wrapped{}, exception.NotAuthorised{Permission: apiPermissions.Permission(jsonRpcMethod)}
}

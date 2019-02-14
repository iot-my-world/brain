package apiAuth

import (
	brainException "gitlab.com/iotTracker/brain/exception"
	apiAuthException "gitlab.com/iotTracker/brain/security/apiAuth/exception"
	"gitlab.com/iotTracker/brain/security/permission"
	permissionHandler "gitlab.com/iotTracker/brain/security/permission/handler"
	"gitlab.com/iotTracker/brain/security/token"
	"gitlab.com/iotTracker/brain/security/claims/login"
	"gitlab.com/iotTracker/brain/security/claims/registerCompanyAdminUser"
	"gitlab.com/iotTracker/brain/security/claims"
)

type APIAuthorizer struct {
	JWTValidator      token.JWTValidator
	PermissionHandler permissionHandler.Handler
}

func (a *APIAuthorizer) AuthorizeAPIReq(jwt string, jsonRpcMethod string) (claims.Claims, error) {

	// Validate the jwt
	jwtClaims, err := a.JWTValidator.ValidateJWT(jwt)
	if err != nil {
		return nil, err
	}

	switch c := jwtClaims.(type) {
	case login.Login:
		// if these are login claims we check in the normal way if the user has the
		// required permission to check access the api
		userHasPermissionResponse := permissionHandler.UserHasPermissionResponse{}
		if err := a.PermissionHandler.UserHasPermission(&permissionHandler.UserHasPermissionRequest{
			UserIdentifier: c.UserId,
			Permission:     permission.Permission(jsonRpcMethod),
		}, &userHasPermissionResponse); err != nil {
			return nil, brainException.Unexpected{Reasons: []string{"determining if user has permission", err.Error()}}
		}

		if !userHasPermissionResponse.Result {
			return nil, apiAuthException.NotAuthorised{Permission: permission.Permission(jsonRpcMethod)}
		}

	case registerCompanyAdminUser.RegisterCompanyAdminUser:
		permissionForMethod := permission.Permission(jsonRpcMethod)
		// check the permissions granted by the RegisterCompanyAdminUser claims to see if this
		// method is allowed
		for allowedPermIdx := range registerCompanyAdminUser.GrantedAPIPermissions {
			if registerCompanyAdminUser.GrantedAPIPermissions[allowedPermIdx] == permissionForMethod {
				return jwtClaims, nil
			}
			if allowedPermIdx == len(registerCompanyAdminUser.GrantedAPIPermissions)-1 {
				return nil, apiAuthException.NotAuthorised{Permission: permission.Permission(jsonRpcMethod)}
			}
		}

	default:
		return nil, apiAuthException.NotAuthorised{Permission: permission.Permission(jsonRpcMethod)}
	}

	return nil, apiAuthException.NotAuthorised{Permission: permission.Permission(jsonRpcMethod)}
}

package apiAuth

import (
	globalException "gitlab.com/iotTracker/brain/exception"
	apiAuthException "gitlab.com/iotTracker/brain/security/apiAuth/exception"
	"gitlab.com/iotTracker/brain/security/permission"
	permissionHandler "gitlab.com/iotTracker/brain/security/permission/handler"
	"gitlab.com/iotTracker/brain/security/token"
	"gitlab.com/iotTracker/brain/security/claims/login"
)

type APIAuthorizer struct {
	JWTValidator      token.JWTValidator
	PermissionHandler permissionHandler.Handler
}

func (a *APIAuthorizer) AuthorizeAPIReq(jwt string, jsonRpcMethod string) error {

	// Validate the jwt
	jwtClaims, err := a.JWTValidator.ValidateJWT(jwt)
	if err != nil {
		return err
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
			return globalException.Unexpected{Reasons: []string{"determining if user has permission", err.Error()}}
		}

		if !userHasPermissionResponse.Result {
			return apiAuthException.NotAuthorised{Permission: permission.Permission(jsonRpcMethod)}
		}
	}


	return nil
}

package apiAuth

import (
	globalException "gitlab.com/iotTracker/brain/exception"
	apiAuthException "gitlab.com/iotTracker/brain/security/apiAuth/exception"
	"gitlab.com/iotTracker/brain/security/permission"
	"gitlab.com/iotTracker/brain/security/token"
)

type APIAuthorizer struct {
	JWTValidator      token.JWTValidator
	PermissionHandler permission.Handler
}

func (a *APIAuthorizer) AuthorizeAPIReq(jwt string, jsonRpcMethod string) error {

	// Validate the jwt
	jwtClaims, err := a.JWTValidator.ValidateJWT(jwt)
	if err != nil {
		return err
	}

	// Check the if the user is authorised to access this jsonRpcMethod based on their role claim
	userHasPermissionResponse := permission.UserHasPermissionResponse{}
	if err := a.PermissionHandler.UserHasPermission(&permission.UserHasPermissionRequest{
		UserIdentifier: jwtClaims.UserId,
		Permission:     permission.Permission(jsonRpcMethod),
	}, &userHasPermissionResponse); err != nil {
		return globalException.Unexpected{Reasons: []string{"determining if user has permission", err.Error()}}
	}

	if !userHasPermissionResponse.Result {
		return apiAuthException.NotAuthorised{Permission: permission.Permission(jsonRpcMethod)}
	}

	return nil
}

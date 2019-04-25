package api

import (
	brainException "gitlab.com/iotTracker/brain/exception"
	apiAuthorizer "gitlab.com/iotTracker/brain/security/authorization/api/authorizer"
	apiAuthException "gitlab.com/iotTracker/brain/security/authorization/api/authorizer/exception"
	apiUserLoginClaims "gitlab.com/iotTracker/brain/security/claims/login/user/api"
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
	case apiUserLoginClaims.Login:
		// if these are login claims we check in the normal way if the user has the
		// required permission to check access the api
		userHasPermissionResponse, err := a.permissionHandler.UserHasPermission(&permissionAdministrator.UserHasPermissionRequest{
			Claims:         typedClaims,
			UserIdentifier: typedClaims.UserId,
			Permission:     api.Permission(jsonRpcMethod),
		})
		if err != nil {
			return wrappedClaims.Wrapped{}, brainException.Unexpected{Reasons: []string{"determining if api user has permission", err.Error()}}
		}
		if !userHasPermissionResponse.Result {
			return wrappedClaims.Wrapped{}, apiAuthException.NotAuthorised{Permission: api.Permission(jsonRpcMethod)}
		}
		// api user was authorised
		return wrappedJWTClaims, nil

	default:
		return wrappedClaims.Wrapped{}, apiAuthException.NotAuthorised{Permission: api.Permission(jsonRpcMethod)}
	}

	return wrappedClaims.Wrapped{}, apiAuthException.NotAuthorised{Permission: api.Permission(jsonRpcMethod)}
}

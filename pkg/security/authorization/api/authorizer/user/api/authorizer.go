package api

import (
	brainException "github.com/iot-my-world/brain/internal/exception"
	authorizer2 "github.com/iot-my-world/brain/pkg/security/authorization/api/authorizer"
	"github.com/iot-my-world/brain/pkg/security/authorization/api/authorizer/exception"
	"github.com/iot-my-world/brain/pkg/security/claims/login/user/api"
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
	case api.Login:
		// if these are login claims we check in the normal way if the user has the
		// required permission to check access the api
		userHasPermissionResponse, err := a.permissionHandler.UserHasPermission(&administrator.UserHasPermissionRequest{
			Claims:         typedClaims,
			UserIdentifier: typedClaims.UserId,
			Permission:     api2.Permission(jsonRpcMethod),
		})
		if err != nil {
			return wrapped.Wrapped{}, brainException.Unexpected{Reasons: []string{"determining if api user has permission", err.Error()}}
		}
		if !userHasPermissionResponse.Result {
			return wrapped.Wrapped{}, exception.NotAuthorised{Permission: api2.Permission(jsonRpcMethod)}
		}
		// api user was authorised
		return wrappedJWTClaims, nil

	default:
		return wrapped.Wrapped{}, exception.NotAuthorised{Permission: api2.Permission(jsonRpcMethod)}
	}

	return wrapped.Wrapped{}, exception.NotAuthorised{Permission: api2.Permission(jsonRpcMethod)}
}

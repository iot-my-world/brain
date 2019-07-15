package authoriser

import (
	jsonRpcServerAuthoriser "github.com/iot-my-world/brain/pkg/api/jsonRpc/server/authoriser"
	authoriserException "github.com/iot-my-world/brain/pkg/api/jsonRpc/server/authoriser/exception"
	"github.com/iot-my-world/brain/pkg/security/claims"
	sigfoxBackendClaims "github.com/iot-my-world/brain/pkg/security/claims/sigfoxBackend"
	wrappedClaims "github.com/iot-my-world/brain/pkg/security/claims/wrapped"
	apiPermissions "github.com/iot-my-world/brain/pkg/security/permission/api"
	"github.com/iot-my-world/brain/pkg/security/token"
)

type authoriser struct {
	jwtValidator token.JWTValidator
}

func New(
	jwtValidator token.JWTValidator,
) jsonRpcServerAuthoriser.Authoriser {
	return &authoriser{
		jwtValidator: jwtValidator,
	}
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

	switch unwrappedJWTClaims.(type) {
	case sigfoxBackendClaims.SigfoxBackend:
		// check the permissions granted by the SigfoxBackendClaims claims to see if this method is allowed
		for allowedPermIdx := range sigfoxBackendClaims.GrantedAPIPermissions {
			permissionForMethod := apiPermissions.Permission(jsonRpcMethod)
			if sigfoxBackendClaims.GrantedAPIPermissions[allowedPermIdx] == permissionForMethod {
				return wrappedJWTClaims, nil
			}
		}

	default:
		return wrappedClaims.Wrapped{}, authoriserException.InvalidClaims{ExpectedClaimsType: claims.SigfoxBackend}
	}
	return wrappedClaims.Wrapped{}, authoriserException.NotAuthorised{Permission: apiPermissions.Permission(jsonRpcMethod)}
}

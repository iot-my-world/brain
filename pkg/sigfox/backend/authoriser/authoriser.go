package authoriser

import (
	jsonRpcServerAuthoriser "github.com/iot-my-world/brain/pkg/api/jsonRpc/server/authoriser"
	wrappedClaims "github.com/iot-my-world/brain/pkg/security/claims/wrapped"
)

type authorizer struct {
}

func New() jsonRpcServerAuthoriser.Authoriser {
	return &authorizer{}
}

func (a *authorizer) AuthoriseServiceMethod(jwt string, jsonRpcMethod string) (wrappedClaims.Wrapped, error) {
	return wrappedClaims.Wrapped{}, nil
}

package authoriser

import (
	wrappedClaims "github.com/iot-my-world/brain/pkg/security/claims/wrapped"
)

type Authoriser interface {
	AuthoriseServiceMethod(jwt string, jsonRpcMethod string) (wrappedClaims.Wrapped, error)
}

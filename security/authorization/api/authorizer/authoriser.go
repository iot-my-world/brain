package authorizer

import (
	wrappedClaims "github.com/iot-my-world/brain/security/claims/wrapped"
)

type Authorizer interface {
	AuthorizeAPIReq(jwt string, jsonRpcMethod string) (wrappedClaims.Wrapped, error)
}

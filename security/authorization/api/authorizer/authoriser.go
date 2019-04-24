package authorizer

import (
	wrappedClaims "gitlab.com/iotTracker/brain/security/claims/wrapped"
)

type Authorizer interface {
	AuthorizeAPIReq(jwt string, jsonRpcMethod string) (wrappedClaims.Wrapped, error)
}

package authorizer

import (
	"github.com/iot-my-world/brain/pkg/security/claims/wrapped"
)

type Authorizer interface {
	AuthorizeAPIReq(jwt string, jsonRpcMethod string) (wrapped.Wrapped, error)
}

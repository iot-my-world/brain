package authenticator

import (
	"gitlab.com/iotTracker/brain/search/identifier"
	"gitlab.com/iotTracker/brain/security/claims"
	"gitlab.com/iotTracker/brain/tracker/zx303"
)

type Authenticator interface {
	Login(request *LoginRequest) (*LoginResponse, error)
}

type LoginRequest struct {
	Claims     claims.Claims
	Identifier identifier.Identifier
}

type LoginResponse struct {
	Result bool
	ZX303  zx303.ZX303
}

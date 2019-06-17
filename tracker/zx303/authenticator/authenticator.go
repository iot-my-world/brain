package authenticator

import (
	"github.com/iot-my-world/brain/search/identifier"
	"github.com/iot-my-world/brain/security/claims"
	"github.com/iot-my-world/brain/tracker/zx303"
)

type Authenticator interface {
	Login(request *LoginRequest) (*LoginResponse, error)
	Logout(request *LogoutRequest) (*LogoutResponse, error)
}

type LoginRequest struct {
	Claims          claims.Claims
	ZX303Identifier identifier.Identifier
}

type LoginResponse struct {
	Result bool
	ZX303  zx303.ZX303
}

type LogoutRequest struct {
	Claims          claims.Claims
	ZX303Identifier identifier.Identifier
}

type LogoutResponse struct {
}

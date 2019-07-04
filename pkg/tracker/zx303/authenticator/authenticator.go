package authenticator

import (
	"github.com/iot-my-world/brain/pkg/search/identifier"
	"github.com/iot-my-world/brain/pkg/security/claims"
	zx3032 "github.com/iot-my-world/brain/pkg/tracker/zx303"
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
	ZX303  zx3032.ZX303
}

type LogoutRequest struct {
	Claims          claims.Claims
	ZX303Identifier identifier.Identifier
}

type LogoutResponse struct {
}

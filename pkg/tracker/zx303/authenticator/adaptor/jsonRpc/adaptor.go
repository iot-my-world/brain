package jsonRpc

import (
	"github.com/iot-my-world/brain/internal/log"
	wrappedIdentifier "github.com/iot-my-world/brain/pkg/search/identifier/wrapped"
	zx3032 "github.com/iot-my-world/brain/pkg/tracker/zx303"
	"github.com/iot-my-world/brain/pkg/tracker/zx303/authenticator"
	wrappedClaims "github.com/iot-my-world/brain/security/claims/wrapped"
	"net/http"
)

type adaptor struct {
	authenticator authenticator.Authenticator
}

func New(
	authenticator authenticator.Authenticator,
) *adaptor {
	return &adaptor{
		authenticator: authenticator,
	}
}

type LoginRequest struct {
	WrappedZX303Identifier wrappedIdentifier.Wrapped `json:"zx303Identifier"`
}

type LoginResponse struct {
	Result bool         `json:"result"`
	ZX303  zx3032.ZX303 `json:"zx303"`
}

func (a *adaptor) Login(r *http.Request, request *LoginRequest, response *LoginResponse) error {
	claims, err := wrappedClaims.UnwrapClaimsFromContext(r)
	if err != nil {
		log.Warn(err.Error())
		return err
	}

	loginResponse, err := a.authenticator.Login(&authenticator.LoginRequest{
		Claims:          claims,
		ZX303Identifier: request.WrappedZX303Identifier.Identifier,
	})
	if err != nil {
		return err
	}

	response.Result = loginResponse.Result
	response.ZX303 = loginResponse.ZX303

	return nil
}

type LogoutRequest struct {
	WrappedZX303Identifier wrappedIdentifier.Wrapped `json:"zx303Identifier"`
}

type LogoutResponse struct {
}

func (a *adaptor) Logout(r *http.Request, request *LogoutRequest, response *LogoutResponse) error {
	claims, err := wrappedClaims.UnwrapClaimsFromContext(r)
	if err != nil {
		log.Warn(err.Error())
		return err
	}

	if _, err := a.authenticator.Logout(&authenticator.LogoutRequest{
		Claims:          claims,
		ZX303Identifier: request.WrappedZX303Identifier.Identifier,
	}); err != nil {
		return err
	}

	return nil
}

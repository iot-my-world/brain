package jsonRpc

import (
	"gitlab.com/iotTracker/brain/log"
	wrappedIdentifier "gitlab.com/iotTracker/brain/search/identifier/wrapped"
	wrappedClaims "gitlab.com/iotTracker/brain/security/claims/wrapped"
	"gitlab.com/iotTracker/brain/tracker/device/zx303"
	zx303DeviceAuthenticator "gitlab.com/iotTracker/brain/tracker/device/zx303/authenticator"
	"net/http"
)

type adaptor struct {
	authenticator zx303DeviceAuthenticator.Authenticator
}

func New(
	authenticator zx303DeviceAuthenticator.Authenticator,
) *adaptor {
	return &adaptor{
		authenticator: authenticator,
	}
}

type LoginRequest struct {
	WrappedIdentifier wrappedIdentifier.Wrapped `json:"identifier"`
}

type LoginResponse struct {
	Result bool        `json:"result"`
	ZX303  zx303.ZX303 `json:"zx303"`
}

func (a *adaptor) Login(r *http.Request, request *LoginRequest, response *LoginResponse) error {
	claims, err := wrappedClaims.UnwrapClaimsFromContext(r)
	if err != nil {
		log.Warn(err.Error())
		return err
	}

	loginResponse, err := a.authenticator.Login(&zx303DeviceAuthenticator.LoginRequest{
		Claims:     claims,
		Identifier: request.WrappedIdentifier.Identifier,
	})
	if err != nil {
		return err
	}

	response.Result = loginResponse.Result
	response.ZX303 = loginResponse.ZX303

	return nil
}

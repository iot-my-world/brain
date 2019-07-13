package jsonRpc

import (
	jsonRpcServerAuthenticator "github.com/iot-my-world/brain/pkg/api/jsonRpc/server/authenticator"
	jsonRpcServiceProvider "github.com/iot-my-world/brain/pkg/api/jsonRpc/service/provider"
	"net/http"
)

type adaptor struct {
	authorizationAdministrator jsonRpcServerAuthenticator.Authenticator
}

func New(authorizationAdministrator jsonRpcServerAuthenticator.Authenticator) *adaptor {
	return &adaptor{
		authorizationAdministrator: authorizationAdministrator,
	}
}

func (a *adaptor) Name() jsonRpcServiceProvider.Name {
	return jsonRpcServiceProvider.Name(jsonRpcServerAuthenticator.ServiceProvider)
}

func (a *adaptor) MethodRequiresAuthorization(method string) bool {
	switch method {
	case jsonRpcServerAuthenticator.LoginService:
		return false
	}
	return true
}

type LogoutRequest struct {
}

type LogoutResponse struct {
}

func (a *adaptor) Logout(r *http.Request, request *LogoutRequest, response *LogoutResponse) error {
	return nil
}

type LoginRequest struct {
	UsernameOrEmailAddress string `json:"usernameOrEmailAddress"`
	Password               string `json:"password"`
}

type LoginResponse struct {
	Jwt string `json:"jwt"`
}

func (a *adaptor) Login(r *http.Request, request *LoginRequest, response *LoginResponse) error {

	loginResponse, err := a.authorizationAdministrator.Login(&jsonRpcServerAuthenticator.LoginRequest{
		UsernameOrEmailAddress: request.UsernameOrEmailAddress,
		Password:               request.Password,
	})
	if err != nil {
		return err
	}

	response.Jwt = loginResponse.Jwt

	return nil
}

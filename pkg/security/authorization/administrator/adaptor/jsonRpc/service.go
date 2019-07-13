package auth

import (
	"fmt"
	jsonRpcServiceProvider "github.com/iot-my-world/brain/pkg/api/jsonRpc/service/provider"
	"github.com/iot-my-world/brain/pkg/security/authorization/administrator"
	"net/http"
)

type adaptor struct {
	authorizationAdministrator administrator.Administrator
}

func New(authorizationAdministrator administrator.Administrator) *adaptor {
	return &adaptor{
		authorizationAdministrator: authorizationAdministrator,
	}
}

func (a *adaptor) Name() jsonRpcServiceProvider.Name {
	return jsonRpcServiceProvider.Name(administrator.ServiceProvider)
}

func (a *adaptor) MethodRequiresAuthorization(method string) bool {
	switch method {
	case administrator.LoginService:
		return false
	}
	return true
}

type LogoutRequest struct {
}

type LogoutResponse struct {
}

func (a *adaptor) Logout(r *http.Request, request *LogoutRequest, response *LogoutResponse) error {
	fmt.Println("Logout Service running.")
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

	loginResponse, err := a.authorizationAdministrator.Login(&administrator.LoginRequest{
		UsernameOrEmailAddress: request.UsernameOrEmailAddress,
		Password:               request.Password,
	})
	if err != nil {
		return err
	}

	response.Jwt = loginResponse.Jwt

	return nil
}

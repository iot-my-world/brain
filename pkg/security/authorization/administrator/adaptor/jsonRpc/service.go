package auth

import (
	"fmt"
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

type LogoutRequest struct {
}

type LogoutResponse struct {
}

func (s *adaptor) Logout(r *http.Request, request *LogoutRequest, response *LogoutResponse) error {
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

func (s *adaptor) Login(r *http.Request, request *LoginRequest, response *LoginResponse) error {

	loginResponse, err := s.authorizationAdministrator.Login(&administrator.LoginRequest{
		UsernameOrEmailAddress: request.UsernameOrEmailAddress,
		Password:               request.Password,
	})
	if err != nil {
		return err
	}

	response.Jwt = loginResponse.Jwt

	return nil
}

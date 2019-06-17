package auth

import (
	"fmt"
	authService "github.com/iot-my-world/brain/security/authorization/service"
	"net/http"
)

type adaptor struct {
	authService authService.Service
}

func New(authService authService.Service) *adaptor {
	return &adaptor{
		authService: authService,
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

	loginResponse, err := s.authService.Login(&authService.LoginRequest{
		UsernameOrEmailAddress: request.UsernameOrEmailAddress,
		Password:               request.Password,
	})
	if err != nil {
		return err
	}

	response.Jwt = loginResponse.Jwt

	return nil
}

package auth

import (
	"fmt"
	"gitlab.com/iotTracker/brain/security/auth"
	"net/http"
)

type adaptor struct {
	authService auth.Service
}

func New(authService auth.Service) *adaptor {
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

	loginRequest := auth.LoginRequest{
		UsernameOrEmailAddress: request.UsernameOrEmailAddress,
		Password:               request.Password,
	}
	loginResponse := auth.LoginResponse{}

	if err := s.authService.Login(&loginRequest, &loginResponse); err != nil {
		return err
	}

	response.Jwt = loginResponse.Jwt

	return nil
}

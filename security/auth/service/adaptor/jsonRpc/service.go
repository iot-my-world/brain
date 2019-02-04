package auth

import (
	"net/http"
	"fmt"
	"gitlab.com/iotTracker/brain/party"
	"gitlab.com/iotTracker/brain/security/auth"
)

type service struct {
	authService auth.Service
}

func New(authService auth.Service) *service {
	return &service{
		authService: authService,
	}
}

type LogoutRequest struct {
}

type LogoutResponse struct {
}

func (s *service) Logout(r *http.Request, request *LogoutRequest, response *LogoutResponse) error {
	fmt.Println("Logout Service running.")
	return nil
}

type LoginRequest struct {
	UsernameOrEmailAddress string `json:"usernameOrEmailAddress"`
	Password               string `json:"password"`
}

type LoginResponse struct {
	Jwt  string     `json:"jwt"`
	User party.User `json:"user"`
}

func (s *service) Login(r *http.Request, request *LoginRequest, response *LoginResponse) error {

	loginRequest := auth.LoginRequest{
		UsernameOrEmailAddress: request.UsernameOrEmailAddress,
		Password:               request.Password,
	}
	loginResponse := auth.LoginResponse{}

	if err := s.authService.Login(&loginRequest, &loginResponse); err != nil {
		return err
	}

	response.Jwt = loginResponse.Jwt
	response.User = loginResponse.User

	return nil
}

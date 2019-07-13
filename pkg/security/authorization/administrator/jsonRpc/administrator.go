package jsonRpc

import (
	jsonRpcClient "github.com/iot-my-world/brain/internal/api/jsonRpc/client"
	administrator2 "github.com/iot-my-world/brain/pkg/security/authorization/administrator"
	"github.com/iot-my-world/brain/pkg/security/authorization/administrator/adaptor/jsonRpc"
)

type administrator struct {
	jsonRpcClient jsonRpcClient.Client
}

func New(
	jsonRpcClient jsonRpcClient.Client,
) administrator2.Administrator {
	return &administrator{
		jsonRpcClient: jsonRpcClient,
	}
}

func (a *administrator) Login(request *administrator2.LoginRequest) (*administrator2.LoginResponse, error) {
	loginResponse := auth.LoginResponse{}
	if err := a.jsonRpcClient.JsonRpcRequest(
		administrator2.LoginService,
		auth.LoginRequest{
			UsernameOrEmailAddress: request.UsernameOrEmailAddress,
			Password:               request.Password,
		},
		&loginResponse); err != nil {
		return nil, err
	}

	return &administrator2.LoginResponse{
		Jwt: loginResponse.Jwt,
	}, nil
}

func (a *administrator) Logout(request *administrator2.LogoutRequest) (*administrator2.LogoutResponse, error) {
	logoutResponse := auth.LogoutResponse{}
	if err := a.jsonRpcClient.JsonRpcRequest(
		administrator2.LogoutService,
		auth.LogoutRequest{},
		&logoutResponse); err != nil {
		return nil, err
	}

	return &administrator2.LogoutResponse{}, nil
}

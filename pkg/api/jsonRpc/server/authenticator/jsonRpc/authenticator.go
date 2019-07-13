package jsonRpc

import (
	"github.com/iot-my-world/brain/internal/log"
	jsonRpcClient "github.com/iot-my-world/brain/pkg/api/jsonRpc/client"
	jsonRpcServerAuthenticator "github.com/iot-my-world/brain/pkg/api/jsonRpc/server/authenticator"
	jsonRpcServerAuthenticatorJsonRpcAdaptor "github.com/iot-my-world/brain/pkg/api/jsonRpc/server/authenticator/adaptor/jsonRpc"
)

type authenticator struct {
	jsonRpcClient jsonRpcClient.Client
}

func New(
	jsonRpcClient jsonRpcClient.Client,
) jsonRpcServerAuthenticator.Authenticator {
	return &authenticator{
		jsonRpcClient: jsonRpcClient,
	}
}

func (a *authenticator) Login(request *jsonRpcServerAuthenticator.LoginRequest) (*jsonRpcServerAuthenticator.LoginResponse, error) {
	loginResponse := jsonRpcServerAuthenticatorJsonRpcAdaptor.LoginResponse{}
	if err := a.jsonRpcClient.JsonRpcRequest(
		jsonRpcServerAuthenticator.LoginService,
		jsonRpcServerAuthenticatorJsonRpcAdaptor.LoginRequest{
			UsernameOrEmailAddress: request.UsernameOrEmailAddress,
			Password:               request.Password,
		},
		&loginResponse,
	); err != nil {
		log.Error(err.Error())
		return nil, err
	}

	return &jsonRpcServerAuthenticator.LoginResponse{Jwt: loginResponse.Jwt}, nil
}

func (a *authenticator) Logout(request *jsonRpcServerAuthenticator.LogoutRequest) (*jsonRpcServerAuthenticator.LogoutResponse, error) {
	logoutResponse := jsonRpcServerAuthenticatorJsonRpcAdaptor.LogoutResponse{}
	if err := a.jsonRpcClient.JsonRpcRequest(
		jsonRpcServerAuthenticator.LogoutService,
		jsonRpcServerAuthenticatorJsonRpcAdaptor.LogoutRequest{},
		&logoutResponse,
	); err != nil {
		log.Error(err.Error())
		return nil, err
	}

	return &jsonRpcServerAuthenticator.LogoutResponse{}, nil
}

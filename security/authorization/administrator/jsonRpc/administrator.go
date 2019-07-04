package jsonRpc

import (
	jsonRpcClient "github.com/iot-my-world/brain/pkg/communication/jsonRpc/client"
	authorizationAdministrator "github.com/iot-my-world/brain/security/authorization/administrator"
	authorizationAdministratorJsonRpcAdaptor "github.com/iot-my-world/brain/security/authorization/administrator/adaptor/jsonRpc"
)

type administrator struct {
	jsonRpcClient jsonRpcClient.Client
}

func New(
	jsonRpcClient jsonRpcClient.Client,
) authorizationAdministrator.Administrator {
	return &administrator{
		jsonRpcClient: jsonRpcClient,
	}
}

func (a *administrator) Login(request *authorizationAdministrator.LoginRequest) (*authorizationAdministrator.LoginResponse, error) {
	loginResponse := authorizationAdministratorJsonRpcAdaptor.LoginResponse{}
	if err := a.jsonRpcClient.JsonRpcRequest(
		authorizationAdministrator.LoginService,
		authorizationAdministratorJsonRpcAdaptor.LoginRequest{
			UsernameOrEmailAddress: request.UsernameOrEmailAddress,
			Password:               request.Password,
		},
		&loginResponse); err != nil {
		return nil, err
	}

	return &authorizationAdministrator.LoginResponse{
		Jwt: loginResponse.Jwt,
	}, nil
}

func (a *administrator) Logout(request *authorizationAdministrator.LogoutRequest) (*authorizationAdministrator.LogoutResponse, error) {
	logoutResponse := authorizationAdministratorJsonRpcAdaptor.LogoutResponse{}
	if err := a.jsonRpcClient.JsonRpcRequest(
		authorizationAdministrator.LogoutService,
		authorizationAdministratorJsonRpcAdaptor.LogoutRequest{},
		&logoutResponse); err != nil {
		return nil, err
	}

	return &authorizationAdministrator.LogoutResponse{}, nil
}

package jsonRpc

import (
	brainException "github.com/iot-my-world/brain/internal/exception"
	jsonRpcClient "github.com/iot-my-world/brain/pkg/communication/jsonRpc/client"
	wrappedIdentifier "github.com/iot-my-world/brain/pkg/search/identifier/wrapped"
	authenticator2 "github.com/iot-my-world/brain/pkg/tracker/zx303/authenticator"
	"github.com/iot-my-world/brain/pkg/tracker/zx303/authenticator/adaptor/jsonRpc"
)

type authenticator struct {
	jsonRpcClient jsonRpcClient.Client
}

func New(
	jsonRpcClient jsonRpcClient.Client,
) authenticator2.Authenticator {
	return &authenticator{
		jsonRpcClient: jsonRpcClient,
	}
}

func (a *authenticator) ValidateLoginRequest(request *authenticator2.LoginRequest) error {
	reasonsInvalid := make([]string, 0)

	if request.ZX303Identifier == nil {
		reasonsInvalid = append(reasonsInvalid, "identifier is nil")
	}

	if !a.jsonRpcClient.LoggedIn() {
		reasonsInvalid = append(reasonsInvalid, "json rpc client is not logged in")
	}

	if len(reasonsInvalid) > 0 {
		return brainException.RequestInvalid{Reasons: reasonsInvalid}
	}
	return nil
}

func (a *authenticator) Login(request *authenticator2.LoginRequest) (*authenticator2.LoginResponse, error) {
	if err := a.ValidateLoginRequest(request); err != nil {
		return nil, err
	}

	// create wrapped identifier
	wrappedDeviceIdentifier, err := wrappedIdentifier.Wrap(request.ZX303Identifier)
	if err != nil {
		return nil, brainException.Unexpected{Reasons: []string{"wrapping device identifier", err.Error()}}
	}

	// login the device
	loginResponse := jsonRpc.LoginResponse{}
	if err := a.jsonRpcClient.JsonRpcRequest(
		"ZX303DeviceAuthenticator.Login",
		jsonRpc.LoginRequest{
			WrappedZX303Identifier: *wrappedDeviceIdentifier,
		},
		&loginResponse,
	); err != nil {
		return nil, err
	}

	return &authenticator2.LoginResponse{
		Result: loginResponse.Result,
		ZX303:  loginResponse.ZX303,
	}, nil
}

func (a *authenticator) ValidateLogoutRequest(request *authenticator2.LogoutRequest) error {
	reasonsInvalid := make([]string, 0)

	if request.ZX303Identifier == nil {
		reasonsInvalid = append(reasonsInvalid, "identifier is nil")
	}

	if !a.jsonRpcClient.LoggedIn() {
		reasonsInvalid = append(reasonsInvalid, "json rpc client is not logged in")
	}

	if len(reasonsInvalid) > 0 {
		return brainException.RequestInvalid{Reasons: reasonsInvalid}
	}
	return nil
}

func (a *authenticator) Logout(request *authenticator2.LogoutRequest) (*authenticator2.LogoutResponse, error) {
	if err := a.ValidateLogoutRequest(request); err != nil {
		return nil, err
	}

	// create wrapped identifier
	wrappedDeviceIdentifier, err := wrappedIdentifier.Wrap(request.ZX303Identifier)
	if err != nil {
		return nil, brainException.Unexpected{Reasons: []string{"wrapping device identifier", err.Error()}}
	}

	// login the device
	loginResponse := jsonRpc.LoginResponse{}
	if err := a.jsonRpcClient.JsonRpcRequest(
		"ZX303DeviceAuthenticator.Logout",
		jsonRpc.LoginRequest{
			WrappedZX303Identifier: *wrappedDeviceIdentifier,
		},
		&loginResponse,
	); err != nil {
		return nil, err
	}

	return &authenticator2.LogoutResponse{}, nil
}

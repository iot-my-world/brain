package jsonRpc

import (
	jsonRpcClient "github.com/iot-my-world/brain/communication/jsonRpc/client"
	brainException "github.com/iot-my-world/brain/exception"
	wrappedIdentifier "github.com/iot-my-world/brain/search/identifier/wrapped"
	zx303DeviceAuthenticator "github.com/iot-my-world/brain/tracker/zx303/authenticator"
	zx303DeviceAuthenticatorJsonRpcAdaptor "github.com/iot-my-world/brain/tracker/zx303/authenticator/adaptor/jsonRpc"
)

type authenticator struct {
	jsonRpcClient jsonRpcClient.Client
}

func New(
	jsonRpcClient jsonRpcClient.Client,
) zx303DeviceAuthenticator.Authenticator {
	return &authenticator{
		jsonRpcClient: jsonRpcClient,
	}
}

func (a *authenticator) ValidateLoginRequest(request *zx303DeviceAuthenticator.LoginRequest) error {
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

func (a *authenticator) Login(request *zx303DeviceAuthenticator.LoginRequest) (*zx303DeviceAuthenticator.LoginResponse, error) {
	if err := a.ValidateLoginRequest(request); err != nil {
		return nil, err
	}

	// create wrapped identifier
	wrappedDeviceIdentifier, err := wrappedIdentifier.Wrap(request.ZX303Identifier)
	if err != nil {
		return nil, brainException.Unexpected{Reasons: []string{"wrapping device identifier", err.Error()}}
	}

	// login the device
	loginResponse := zx303DeviceAuthenticatorJsonRpcAdaptor.LoginResponse{}
	if err := a.jsonRpcClient.JsonRpcRequest(
		"ZX303DeviceAuthenticator.Login",
		zx303DeviceAuthenticatorJsonRpcAdaptor.LoginRequest{
			WrappedZX303Identifier: *wrappedDeviceIdentifier,
		},
		&loginResponse,
	); err != nil {
		return nil, err
	}

	return &zx303DeviceAuthenticator.LoginResponse{
		Result: loginResponse.Result,
		ZX303:  loginResponse.ZX303,
	}, nil
}

func (a *authenticator) ValidateLogoutRequest(request *zx303DeviceAuthenticator.LogoutRequest) error {
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

func (a *authenticator) Logout(request *zx303DeviceAuthenticator.LogoutRequest) (*zx303DeviceAuthenticator.LogoutResponse, error) {
	if err := a.ValidateLogoutRequest(request); err != nil {
		return nil, err
	}

	// create wrapped identifier
	wrappedDeviceIdentifier, err := wrappedIdentifier.Wrap(request.ZX303Identifier)
	if err != nil {
		return nil, brainException.Unexpected{Reasons: []string{"wrapping device identifier", err.Error()}}
	}

	// login the device
	loginResponse := zx303DeviceAuthenticatorJsonRpcAdaptor.LoginResponse{}
	if err := a.jsonRpcClient.JsonRpcRequest(
		"ZX303DeviceAuthenticator.Logout",
		zx303DeviceAuthenticatorJsonRpcAdaptor.LoginRequest{
			WrappedZX303Identifier: *wrappedDeviceIdentifier,
		},
		&loginResponse,
	); err != nil {
		return nil, err
	}

	return &zx303DeviceAuthenticator.LogoutResponse{}, nil
}

package jsonRpc

import (
	jsonRpcClient "gitlab.com/iotTracker/brain/communication/jsonRpc/client"
	brainException "gitlab.com/iotTracker/brain/exception"
	wrappedIdentifier "gitlab.com/iotTracker/brain/search/identifier/wrapped"
	zx303DeviceAuthenticator "gitlab.com/iotTracker/brain/tracker/device/zx303/authenticator"
	zx303DeviceAuthenticatorJsonRpcAdaptor "gitlab.com/iotTracker/brain/tracker/device/zx303/authenticator/adaptor/jsonRpc"
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

	if request.Identifier == nil {
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

	// create wrapped identifier
	wrappedDeviceIdentifier, err := wrappedIdentifier.Wrap(request.Identifier)
	if err != nil {
		return nil, brainException.Unexpected{Reasons: []string{"wrapping device identifier", err.Error()}}
	}

	// login the device
	loginResponse := zx303DeviceAuthenticatorJsonRpcAdaptor.LoginResponse{}
	if err := a.jsonRpcClient.JsonRpcRequest(
		"ZX303DeviceAuthenticator.Login",
		zx303DeviceAuthenticatorJsonRpcAdaptor.LoginRequest{
			WrappedIdentifier: *wrappedDeviceIdentifier,
		},
		&loginResponse,
	); err != nil {
		return nil, brainException.Unexpected{Reasons: []string{"log in error", err.Error()}}
	}

	return &zx303DeviceAuthenticator.LoginResponse{
		Result: loginResponse.Result,
		ZX303:  loginResponse.ZX303,
	}, nil
}

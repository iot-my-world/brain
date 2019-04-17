package basic

import (
	brainException "gitlab.com/iotTracker/brain/exception"
	"gitlab.com/iotTracker/brain/tracker/device/zx303"
	zx303DeviceAuthenticator "gitlab.com/iotTracker/brain/tracker/device/zx303/authenticator"
	zx303DeviceAuthenticatorException "gitlab.com/iotTracker/brain/tracker/device/zx303/authenticator/exception"
	zx303RecordHandler "gitlab.com/iotTracker/brain/tracker/device/zx303/recordHandler"
	zx303RecordHandlerException "gitlab.com/iotTracker/brain/tracker/device/zx303/recordHandler/exception"
)

type authenticator struct {
	zx303RecordHandler *zx303RecordHandler.RecordHandler
}

func New(
	zx303RecordHandler *zx303RecordHandler.RecordHandler,
) zx303DeviceAuthenticator.Authenticator {
	return &authenticator{
		zx303RecordHandler: zx303RecordHandler,
	}
}

func (a *authenticator) ValidateLoginRequest(request *zx303DeviceAuthenticator.LoginRequest) error {
	reasonsInvalid := make([]string, 0)

	if request.Claims == nil {
		reasonsInvalid = append(reasonsInvalid, "claims are nil")
	}
	if request.Identifier == nil {
		reasonsInvalid = append(reasonsInvalid, "identifier is nil")
	} else if !zx303.IsValidIdentifier(request.Identifier) {
		reasonsInvalid = append(reasonsInvalid, "identifier is not valid")
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

	// try and retrieve the device
	retrieveResponse, err := a.zx303RecordHandler.Retrieve(&zx303RecordHandler.RetrieveRequest{
		Claims:     request.Claims,
		Identifier: request.Identifier,
	})
	if err != nil {
		switch err.(type) {
		case zx303RecordHandlerException.NotFound:
			// device not found, login not successful
			return &zx303DeviceAuthenticator.LoginResponse{
				Result: false,
			}, nil
		default:
			// some other error retrieving the device
			return nil, zx303DeviceAuthenticatorException.Retrieval{Reasons: []string{"device retrieval", err.Error()}}
		}
	}

	return &zx303DeviceAuthenticator.LoginResponse{
		Result: true,
		ZX303:  retrieveResponse.ZX303,
	}, nil
}

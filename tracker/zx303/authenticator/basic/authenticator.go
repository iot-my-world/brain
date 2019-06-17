package basic

import (
	brainException "github.com/iot-my-world/brain/exception"
	"github.com/iot-my-world/brain/log"
	"github.com/iot-my-world/brain/search/identifier/id"
	"github.com/iot-my-world/brain/tracker/zx303"
	zx303DeviceAuthenticator "github.com/iot-my-world/brain/tracker/zx303/authenticator"
	zx303DeviceAuthenticatorException "github.com/iot-my-world/brain/tracker/zx303/authenticator/exception"
	zx303RecordHandler "github.com/iot-my-world/brain/tracker/zx303/recordHandler"
	zx303RecordHandlerException "github.com/iot-my-world/brain/tracker/zx303/recordHandler/exception"
	"time"
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
	if request.ZX303Identifier == nil {
		reasonsInvalid = append(reasonsInvalid, "identifier is nil")
	} else if !zx303.IsValidIdentifier(request.ZX303Identifier) {
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
		Identifier: request.ZX303Identifier,
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
			err = zx303DeviceAuthenticatorException.Login{Reasons: []string{"device retrieval", err.Error()}}
			log.Error(err.Error())
			return nil, err
		}
	}

	// update the device to logged in
	retrieveResponse.ZX303.LoggedIn = true
	retrieveResponse.ZX303.LogInTimestamp = time.Now().Unix()
	if _, err := a.zx303RecordHandler.Update(&zx303RecordHandler.UpdateRequest{
		Claims: request.Claims,
		Identifier: id.Identifier{
			Id: retrieveResponse.ZX303.Id,
		},
		ZX303: retrieveResponse.ZX303,
	}); err != nil {
		err = zx303DeviceAuthenticatorException.Login{Reasons: []string{"device update", err.Error()}}
		log.Error(err.Error())
		return nil, err
	}

	return &zx303DeviceAuthenticator.LoginResponse{
		Result: true,
		ZX303:  retrieveResponse.ZX303,
	}, nil
}

func (a *authenticator) ValidateLogoutRequest(request *zx303DeviceAuthenticator.LogoutRequest) error {
	reasonsInvalid := make([]string, 0)

	if request.Claims == nil {
		reasonsInvalid = append(reasonsInvalid, "claims are nil")
	}
	if request.ZX303Identifier == nil {
		reasonsInvalid = append(reasonsInvalid, "identifier is nil")
	} else if !zx303.IsValidIdentifier(request.ZX303Identifier) {
		reasonsInvalid = append(reasonsInvalid, "identifier is not valid")
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

	// try and retrieve the device
	retrieveResponse, err := a.zx303RecordHandler.Retrieve(&zx303RecordHandler.RetrieveRequest{
		Claims:     request.Claims,
		Identifier: request.ZX303Identifier,
	})
	if err != nil {
		err = zx303DeviceAuthenticatorException.Logout{Reasons: []string{"device retrieval", err.Error()}}
		log.Error(err.Error())
		return nil, err
	}

	// update the device to logged out
	retrieveResponse.ZX303.LoggedIn = false
	retrieveResponse.ZX303.LogOutTimestamp = time.Now().Unix()
	if _, err := a.zx303RecordHandler.Update(&zx303RecordHandler.UpdateRequest{
		Claims: request.Claims,
		Identifier: id.Identifier{
			Id: retrieveResponse.ZX303.Id,
		},
		ZX303: retrieveResponse.ZX303,
	}); err != nil {
		err = zx303DeviceAuthenticatorException.Logout{Reasons: []string{"device update", err.Error()}}
		log.Error(err.Error())
		return nil, err
	}

	return &zx303DeviceAuthenticator.LogoutResponse{}, nil
}

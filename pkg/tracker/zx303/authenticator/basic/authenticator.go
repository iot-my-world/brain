package basic

import (
	brainException "github.com/iot-my-world/brain/internal/exception"
	"github.com/iot-my-world/brain/internal/log"
	"github.com/iot-my-world/brain/pkg/search/identifier/id"
	zx3032 "github.com/iot-my-world/brain/pkg/tracker/zx303"
	authenticator2 "github.com/iot-my-world/brain/pkg/tracker/zx303/authenticator"
	exception2 "github.com/iot-my-world/brain/pkg/tracker/zx303/authenticator/exception"
	"github.com/iot-my-world/brain/pkg/tracker/zx303/recordHandler"
	"github.com/iot-my-world/brain/pkg/tracker/zx303/recordHandler/exception"
	"time"
)

type authenticator struct {
	zx303RecordHandler *recordHandler.RecordHandler
}

func New(
	zx303RecordHandler *recordHandler.RecordHandler,
) authenticator2.Authenticator {
	return &authenticator{
		zx303RecordHandler: zx303RecordHandler,
	}
}

func (a *authenticator) ValidateLoginRequest(request *authenticator2.LoginRequest) error {
	reasonsInvalid := make([]string, 0)

	if request.Claims == nil {
		reasonsInvalid = append(reasonsInvalid, "claims are nil")
	}
	if request.ZX303Identifier == nil {
		reasonsInvalid = append(reasonsInvalid, "identifier is nil")
	} else if !zx3032.IsValidIdentifier(request.ZX303Identifier) {
		reasonsInvalid = append(reasonsInvalid, "identifier is not valid")
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

	// try and retrieve the device
	retrieveResponse, err := a.zx303RecordHandler.Retrieve(&recordHandler.RetrieveRequest{
		Claims:     request.Claims,
		Identifier: request.ZX303Identifier,
	})
	if err != nil {
		switch err.(type) {
		case exception.NotFound:
			// device not found, login not successful
			return &authenticator2.LoginResponse{
				Result: false,
			}, nil
		default:
			// some other error retrieving the device
			err = exception2.Login{Reasons: []string{"device retrieval", err.Error()}}
			log.Error(err.Error())
			return nil, err
		}
	}

	// update the device to logged in
	retrieveResponse.ZX303.LoggedIn = true
	retrieveResponse.ZX303.LogInTimestamp = time.Now().Unix()
	if _, err := a.zx303RecordHandler.Update(&recordHandler.UpdateRequest{
		Claims: request.Claims,
		Identifier: id.Identifier{
			Id: retrieveResponse.ZX303.Id,
		},
		ZX303: retrieveResponse.ZX303,
	}); err != nil {
		err = exception2.Login{Reasons: []string{"device update", err.Error()}}
		log.Error(err.Error())
		return nil, err
	}

	return &authenticator2.LoginResponse{
		Result: true,
		ZX303:  retrieveResponse.ZX303,
	}, nil
}

func (a *authenticator) ValidateLogoutRequest(request *authenticator2.LogoutRequest) error {
	reasonsInvalid := make([]string, 0)

	if request.Claims == nil {
		reasonsInvalid = append(reasonsInvalid, "claims are nil")
	}
	if request.ZX303Identifier == nil {
		reasonsInvalid = append(reasonsInvalid, "identifier is nil")
	} else if !zx3032.IsValidIdentifier(request.ZX303Identifier) {
		reasonsInvalid = append(reasonsInvalid, "identifier is not valid")
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

	// try and retrieve the device
	retrieveResponse, err := a.zx303RecordHandler.Retrieve(&recordHandler.RetrieveRequest{
		Claims:     request.Claims,
		Identifier: request.ZX303Identifier,
	})
	if err != nil {
		err = exception2.Logout{Reasons: []string{"device retrieval", err.Error()}}
		log.Error(err.Error())
		return nil, err
	}

	// update the device to logged out
	retrieveResponse.ZX303.LoggedIn = false
	retrieveResponse.ZX303.LogOutTimestamp = time.Now().Unix()
	if _, err := a.zx303RecordHandler.Update(&recordHandler.UpdateRequest{
		Claims: request.Claims,
		Identifier: id.Identifier{
			Id: retrieveResponse.ZX303.Id,
		},
		ZX303: retrieveResponse.ZX303,
	}); err != nil {
		err = exception2.Logout{Reasons: []string{"device update", err.Error()}}
		log.Error(err.Error())
		return nil, err
	}

	return &authenticator2.LogoutResponse{}, nil
}

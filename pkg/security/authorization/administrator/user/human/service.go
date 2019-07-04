package human

import (
	"crypto/rsa"
	"errors"
	"fmt"
	"github.com/iot-my-world/brain/pkg/search/identifier/emailAddress"
	"github.com/iot-my-world/brain/pkg/search/identifier/id"
	"github.com/iot-my-world/brain/pkg/search/identifier/username"
	administrator2 "github.com/iot-my-world/brain/pkg/security/authorization/administrator"
	"github.com/iot-my-world/brain/pkg/security/claims/login/user/human"
	token2 "github.com/iot-my-world/brain/pkg/security/token"
	userRecordHandler "github.com/iot-my-world/brain/pkg/user/human/recordHandler"
	userRecordHandlerException "github.com/iot-my-world/brain/pkg/user/human/recordHandler/exception"
	"golang.org/x/crypto/bcrypt"
	"time"
)

type administrator struct {
	userRecordHandler userRecordHandler.RecordHandler
	jwtGenerator      token2.JWTGenerator
	systemClaims      *human.Login
}

func New(
	userRecordHandler userRecordHandler.RecordHandler,
	rsaPrivateKey *rsa.PrivateKey,
	systemClaims *human.Login,
) administrator2.Administrator {
	return &administrator{
		userRecordHandler: userRecordHandler,
		jwtGenerator:      token2.NewJWTGenerator(rsaPrivateKey),
		systemClaims:      systemClaims,
	}
}

func (a *administrator) Logout(request *administrator2.LogoutRequest) (*administrator2.LogoutResponse, error) {
	fmt.Println("Logout Service running.")
	return &administrator2.LogoutResponse{}, nil
}

func (a *administrator) Login(request *administrator2.LoginRequest) (*administrator2.LoginResponse, error) {
	var retrieveUserResponse *userRecordHandler.RetrieveResponse
	var err error

	//try and retrieve User record with username
	retrieveUserResponse, err = a.userRecordHandler.Retrieve(&userRecordHandler.RetrieveRequest{
		Claims:     *a.systemClaims,
		Identifier: username.Identifier{Username: request.UsernameOrEmailAddress},
	})
	if err != nil {
		switch err.(type) {
		case userRecordHandlerException.NotFound:
			//try and retrieve User record with email address
			retrieveUserResponse, err = a.userRecordHandler.Retrieve(&userRecordHandler.RetrieveRequest{
				Claims:     *a.systemClaims,
				Identifier: emailAddress.Identifier{EmailAddress: request.UsernameOrEmailAddress},
			})
			if err != nil {
				return nil, errors.New("log in failed")
			}
		default:
			return nil, errors.New("log in failed")
		}
	}

	//User record retrieved successfully, check password
	if err := bcrypt.CompareHashAndPassword(retrieveUserResponse.User.Password, []byte(request.Password)); err != nil {
		//Password Incorrect
		return nil, errors.New("log In failed")
	}

	// Password is correct. Try and generate loginToken
	loginToken, err := a.jwtGenerator.GenerateToken(human.Login{
		UserId:          id.Identifier{Id: retrieveUserResponse.User.Id},
		IssueTime:       time.Now().UTC().Unix(),
		ExpirationTime:  time.Now().Add(90 * time.Minute).UTC().Unix(),
		ParentPartyType: retrieveUserResponse.User.ParentPartyType,
		ParentId:        retrieveUserResponse.User.ParentId,
		PartyType:       retrieveUserResponse.User.PartyType,
		PartyId:         retrieveUserResponse.User.PartyId,
	})
	if err != nil {
		//Unexpected Error!
		return nil, errors.New("log In failed")
	}

	//Login Successful, return Token
	return &administrator2.LoginResponse{Jwt: loginToken}, nil
}

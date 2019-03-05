package auth

import (
	"crypto/rsa"
	"errors"
	"fmt"
	userRecordHandler "gitlab.com/iotTracker/brain/party/user/recordHandler"
	userRecordHandlerException "gitlab.com/iotTracker/brain/party/user/recordHandler/exception"
	"gitlab.com/iotTracker/brain/search/identifier/emailAddress"
	"gitlab.com/iotTracker/brain/search/identifier/id"
	"gitlab.com/iotTracker/brain/search/identifier/username"
	"gitlab.com/iotTracker/brain/security/auth"
	"gitlab.com/iotTracker/brain/security/claims/login"
	"gitlab.com/iotTracker/brain/security/token"
	"golang.org/x/crypto/bcrypt"
	"time"
)

type service struct {
	userRecordHandler userRecordHandler.RecordHandler
	jwtGenerator      token.JWTGenerator
}

func New(userRecordHandler userRecordHandler.RecordHandler, rsaPrivateKey *rsa.PrivateKey) *service {
	return &service{
		userRecordHandler: userRecordHandler,
		jwtGenerator:      token.NewJWTGenerator(rsaPrivateKey),
	}
}

func (s *service) Logout(request *auth.LogoutRequest, response *auth.LogoutResponse) error {
	fmt.Println("Logout Service running.")
	return nil
}

func (s *service) Login(request *auth.LoginRequest, response *auth.LoginResponse) error {

	retrieveUserResponse := userRecordHandler.RetrieveResponse{}

	//try and retrieve User record with username
	if err := s.userRecordHandler.Retrieve(&userRecordHandler.RetrieveRequest{
		Identifier: username.Identifier{Username: request.UsernameOrEmailAddress},
	}, &retrieveUserResponse); err != nil {
		switch err.(type) {
		case userRecordHandlerException.NotFound:
			//try and retrieve User record with email address
			if err := s.userRecordHandler.Retrieve(&userRecordHandler.RetrieveRequest{
				Identifier: emailAddress.Identifier{EmailAddress: request.UsernameOrEmailAddress},
			}, &retrieveUserResponse); err != nil {
				return errors.New("log in failed")
			}
		default:
			return errors.New("log in failed")
		}
	}

	//User record retrieved successfully, check password
	if err := bcrypt.CompareHashAndPassword(retrieveUserResponse.User.Password, []byte(request.Password)); err != nil {
		//Password Incorrect
		return errors.New("log In failed")
	}

	// Password is correct. Try and generate loginToken
	loginToken, err := s.jwtGenerator.GenerateToken(login.Login{
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
		return errors.New("log In failed")
	}

	//Login Successful, return Token to front-end client
	response.Jwt = loginToken

	return nil
}

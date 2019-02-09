package auth

import (
	"crypto/rsa"
	"errors"
	"fmt"
	"gitlab.com/iotTracker/brain/party/user"
	userException "gitlab.com/iotTracker/brain/party/user/exception"
	"gitlab.com/iotTracker/brain/search/identifier/emailAddress"
	"gitlab.com/iotTracker/brain/search/identifier/id"
	"gitlab.com/iotTracker/brain/search/identifier/username"
	"gitlab.com/iotTracker/brain/security/auth"
	"gitlab.com/iotTracker/brain/security/claims"
	"gitlab.com/iotTracker/brain/security/token"
	"golang.org/x/crypto/bcrypt"
	"time"
)

type service struct {
	userRecordHandler user.RecordHandler
	jwtGenerator      token.JWTGenerator
}

func New(userRecordHandler user.RecordHandler, rsaPrivateKey *rsa.PrivateKey) *service {
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

	retrieveUserResponse := user.RetrieveResponse{}

	//try and retrieve User record with username
	if err := s.userRecordHandler.Retrieve(&user.RetrieveRequest{
		Identifier: username.Identifier{Username: request.UsernameOrEmailAddress},
	}, &retrieveUserResponse); err != nil {
		switch err.(type) {
		case userException.NotFound:
			//try and retrieve User record with email address
			if err := s.userRecordHandler.Retrieve(&user.RetrieveRequest{
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
	loginToken, err := s.jwtGenerator.GenerateLoginToken(claims.Claims{
		UserId:         id.Identifier{Id: retrieveUserResponse.User.Id},
		IssueTime:      time.Now().UTC().Unix(),
		ExpirationTime: time.Now().Add(claims.ValidTime).UTC().Unix(),
		PartyType:      retrieveUserResponse.User.PartyType,
		PartyId:        retrieveUserResponse.User.PartyId,
	})
	if err != nil {
		//Unexpected Error!
		return errors.New("log In failed")
	}

	//Login Successful, return Token to front-end client
	response.Jwt = loginToken

	return nil
}

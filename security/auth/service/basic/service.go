package auth

import (
	"gitlab.com/iotTracker/brain/party/user"
	"golang.org/x/crypto/bcrypt"
	"crypto/rsa"
	"fmt"
	"gitlab.com/iotTracker/brain/security/token"
	"gitlab.com/iotTracker/brain/security/claims"
	"errors"
	"gitlab.com/iotTracker/brain/search/identifier/username"
	"gitlab.com/iotTracker/brain/search/identifier/id"
	userException "gitlab.com/iotTracker/brain/party/user/exception"
	"gitlab.com/iotTracker/brain/search/identifier/emailAddress"
	"gitlab.com/iotTracker/brain/security/auth"
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
	testExpiration := time.Now().AddDate(0, 0, -1).UTC().Unix()
	fmt.Println("test expiration:", testExpiration)
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
	response.User = retrieveUserResponse.User

	return nil
}

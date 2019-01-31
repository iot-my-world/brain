package auth

import (
	"net/http"
	"gitlab.com/iotTracker/brain/party/user"
	"golang.org/x/crypto/bcrypt"
	"crypto/rsa"
	"fmt"
	"gitlab.com/iotTracker/brain/security/token"
	"gitlab.com/iotTracker/brain/security/claims"
	"errors"
	"gitlab.com/iotTracker/brain/search/identifiers/username"
	"gitlab.com/iotTracker/brain/party"
	"gitlab.com/iotTracker/brain/search/identifiers/id"
	userException "gitlab.com/iotTracker/brain/party/user/exception"
	"gitlab.com/iotTracker/brain/search/identifiers/emailAddress"
)

type service struct {
	userRecordHandler user.RecordHandler
	jwtGenerator      token.JWTGenerator
}

func NewService(userRecordHandler user.RecordHandler, rsaPrivateKey *rsa.PrivateKey) *service {
	return &service{
		userRecordHandler: userRecordHandler,
		jwtGenerator:      token.NewJWTGenerator(rsaPrivateKey),
	}
}

type LogoutRequest struct {
}

type LogoutResponse struct {
}

func (s *service) Logout(r *http.Request, request *LogoutRequest, response *LogoutResponse) error {
	fmt.Println("Logout Service running.")
	return nil
}

type LoginRequest struct {
	UsernameOrEmailAddress string `json:"usernameOrEmailAddress"`
	Password               string `json:"password"`
}

type LoginResponse struct {
	Jwt  string     `json:"jwt"`
	User party.User `json:"user"`
}

func (s *service) Login(r *http.Request, request *LoginRequest, response *LoginResponse) error {

	retrieveUserResponse := user.RetrieveResponse{}

	//try and retrieve User record with username
	if err := s.userRecordHandler.Retrieve(&user.RetrieveRequest{
		Identifier: username.Identifier(request.UsernameOrEmailAddress),
	}, &retrieveUserResponse); err != nil {
		switch err.(type) {
		case userException.NotFound:
			//try and retrieve User record with email address
			if err := s.userRecordHandler.Retrieve(&user.RetrieveRequest{
				Identifier: emailAddress.Identifier(request.UsernameOrEmailAddress),
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
	loginToken, err := s.jwtGenerator.GenerateLoginToken(claims.LoginClaims{
		UserId: id.Identifier(retrieveUserResponse.User.Id),
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

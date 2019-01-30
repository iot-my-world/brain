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
	Username string `json:"username" bson:"username"`
	Password string `json:"password" bson:"password"`
}

type LoginResponse struct {
	Jwt  string     `json:"jwt" bson:"jwt"`
	User party.User `json:"user" bson:"user"`
}

func (s *service) Login(r *http.Request, request *LoginRequest, response *LoginResponse) error {

	retrieveUserRequest := user.RetrieveRequest{Identifier: username.Identifier(request.Username)}
	retrieveUserResponse := user.RetrieveResponse{}

	//Retrieve User record
	if err := s.userRecordHandler.Retrieve(&retrieveUserRequest, &retrieveUserResponse); err != nil {
		//Error while retrieving user record
		return errors.New("log In failed")
	}

	//User record retrieved successfully, check password
	if err := bcrypt.CompareHashAndPassword(retrieveUserResponse.User.Password, []byte(request.Password)); err != nil {
		//Password Incorrect
		return errors.New("log In failed")
	}

	//Password is correct. Try and retrieve loginToken
	loginToken, err := s.jwtGenerator.GenerateLoginToken(claims.LoginClaims{
		Username:   request.Username,
		SystemRole: retrieveUserResponse.User.SystemRole,
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

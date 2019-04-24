package ui

import (
	"crypto/rsa"
	"errors"
	"fmt"
	"gitlab.com/iotTracker/brain/search/identifier/emailAddress"
	"gitlab.com/iotTracker/brain/search/identifier/id"
	"gitlab.com/iotTracker/brain/search/identifier/username"
	authService "gitlab.com/iotTracker/brain/security/authorization/service"
	"gitlab.com/iotTracker/brain/security/claims/login"
	"gitlab.com/iotTracker/brain/security/token"
	userRecordHandler "gitlab.com/iotTracker/brain/user/human/recordHandler"
	userRecordHandlerException "gitlab.com/iotTracker/brain/user/human/recordHandler/exception"
	"golang.org/x/crypto/bcrypt"
	"time"
)

type service struct {
	userRecordHandler userRecordHandler.RecordHandler
	jwtGenerator      token.JWTGenerator
	systemClaims      *login.Login
}

func New(
	userRecordHandler userRecordHandler.RecordHandler,
	rsaPrivateKey *rsa.PrivateKey,
	systemClaims *login.Login,
) *service {
	return &service{
		userRecordHandler: userRecordHandler,
		jwtGenerator:      token.NewJWTGenerator(rsaPrivateKey),
		systemClaims:      systemClaims,
	}
}

func (s *service) Logout(request *authService.LogoutRequest) (*authService.LogoutResponse, error) {
	fmt.Println("Logout Service running.")
	return &authService.LogoutResponse{}, nil
}

func (s *service) Login(request *authService.LoginRequest) (*authService.LoginResponse, error) {
	var retrieveUserResponse *userRecordHandler.RetrieveResponse
	var err error

	//try and retrieve User record with username
	retrieveUserResponse, err = s.userRecordHandler.Retrieve(&userRecordHandler.RetrieveRequest{
		Claims:     *s.systemClaims,
		Identifier: username.Identifier{Username: request.UsernameOrEmailAddress},
	})
	if err != nil {
		switch err.(type) {
		case userRecordHandlerException.NotFound:
			//try and retrieve User record with email address
			retrieveUserResponse, err = s.userRecordHandler.Retrieve(&userRecordHandler.RetrieveRequest{
				Claims:     *s.systemClaims,
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
		return nil, errors.New("log In failed")
	}

	//Login Successful, return Token
	return &authService.LoginResponse{Jwt: loginToken}, nil
}

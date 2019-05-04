package human

import (
	"crypto/rsa"
	"errors"
	"fmt"
	"gitlab.com/iotTracker/brain/search/identifier/emailAddress"
	"gitlab.com/iotTracker/brain/search/identifier/id"
	"gitlab.com/iotTracker/brain/search/identifier/username"
	authService "gitlab.com/iotTracker/brain/security/authorization/service"
	apiUserLoginClaims "gitlab.com/iotTracker/brain/security/claims/login/user/api"
	humanUserLoginClaims "gitlab.com/iotTracker/brain/security/claims/login/user/human"
	"gitlab.com/iotTracker/brain/security/token"
	apiUserRecordHandler "gitlab.com/iotTracker/brain/user/api/recordHandler"
	userRecordHandlerException "gitlab.com/iotTracker/brain/user/human/recordHandler/exception"
	"golang.org/x/crypto/bcrypt"
	"time"
)

type service struct {
	apiUserRecordHandler *apiUserRecordHandler.RecordHandler
	jwtGenerator         token.JWTGenerator
	systemClaims         *humanUserLoginClaims.Login
}

func New(
	apiUserRecordHandler *apiUserRecordHandler.RecordHandler,
	rsaPrivateKey *rsa.PrivateKey,
	systemClaims *humanUserLoginClaims.Login,
) *service {
	return &service{
		apiUserRecordHandler: apiUserRecordHandler,
		jwtGenerator:         token.NewJWTGenerator(rsaPrivateKey),
		systemClaims:         systemClaims,
	}
}

func (s *service) Logout(request *authService.LogoutRequest) (*authService.LogoutResponse, error) {
	fmt.Println("Logout Service running.")
	return &authService.LogoutResponse{}, nil
}

func (s *service) Login(request *authService.LoginRequest) (*authService.LoginResponse, error) {
	var retrieveUserResponse *apiUserRecordHandler.RetrieveResponse
	var err error

	//try and retrieve api user record with username
	retrieveUserResponse, err = s.apiUserRecordHandler.Retrieve(&apiUserRecordHandler.RetrieveRequest{
		Claims:     *s.systemClaims,
		Identifier: username.Identifier{Username: request.UsernameOrEmailAddress},
	})
	if err != nil {
		switch err.(type) {
		case userRecordHandlerException.NotFound:
			//try and retrieve User record with email address
			retrieveUserResponse, err = s.apiUserRecordHandler.Retrieve(&apiUserRecordHandler.RetrieveRequest{
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
	loginToken, err := s.jwtGenerator.GenerateToken(apiUserLoginClaims.Login{
		UserId:          id.Identifier{Id: retrieveUserResponse.User.Id},
		IssueTime:       time.Now().UTC().Unix(),
		ExpirationTime:  time.Now().Add(90 * time.Minute).UTC().Unix(),
		ParentPartyType: retrieveUserResponse.User.PartyType,
		ParentId:        retrieveUserResponse.User.PartyId,
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

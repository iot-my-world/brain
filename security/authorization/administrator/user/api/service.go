package human

import (
	"crypto/rsa"
	"errors"
	"github.com/iot-my-world/brain/pkg/search/identifier/emailAddress"
	"github.com/iot-my-world/brain/pkg/search/identifier/id"
	"github.com/iot-my-world/brain/pkg/search/identifier/username"
	apiUserRecordHandler "github.com/iot-my-world/brain/pkg/user/api/recordHandler"
	userRecordHandlerException "github.com/iot-my-world/brain/pkg/user/human/recordHandler/exception"
	authorizationAdministrator "github.com/iot-my-world/brain/security/authorization/administrator"
	apiUserLoginClaims "github.com/iot-my-world/brain/security/claims/login/user/api"
	humanUserLoginClaims "github.com/iot-my-world/brain/security/claims/login/user/human"
	"github.com/iot-my-world/brain/security/token"
	"golang.org/x/crypto/bcrypt"
	"time"
)

type service struct {
	apiUserRecordHandler apiUserRecordHandler.RecordHandler
	jwtGenerator         token.JWTGenerator
	systemClaims         *humanUserLoginClaims.Login
}

func New(
	apiUserRecordHandler apiUserRecordHandler.RecordHandler,
	rsaPrivateKey *rsa.PrivateKey,
	systemClaims *humanUserLoginClaims.Login,
) authorizationAdministrator.Administrator {
	return &service{
		apiUserRecordHandler: apiUserRecordHandler,
		jwtGenerator:         token.NewJWTGenerator(rsaPrivateKey),
		systemClaims:         systemClaims,
	}
}

func (a *service) Logout(request *authorizationAdministrator.LogoutRequest) (*authorizationAdministrator.LogoutResponse, error) {
	return &authorizationAdministrator.LogoutResponse{}, nil
}

func (a *service) Login(request *authorizationAdministrator.LoginRequest) (*authorizationAdministrator.LoginResponse, error) {
	var retrieveUserResponse *apiUserRecordHandler.RetrieveResponse
	var err error

	//try and retrieve api user record with username
	retrieveUserResponse, err = a.apiUserRecordHandler.Retrieve(&apiUserRecordHandler.RetrieveRequest{
		Claims:     *a.systemClaims,
		Identifier: username.Identifier{Username: request.UsernameOrEmailAddress},
	})
	if err != nil {
		switch err.(type) {
		case userRecordHandlerException.NotFound:
			//try and retrieve User record with email address
			retrieveUserResponse, err = a.apiUserRecordHandler.Retrieve(&apiUserRecordHandler.RetrieveRequest{
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
	loginToken, err := a.jwtGenerator.GenerateToken(apiUserLoginClaims.Login{
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
	return &authorizationAdministrator.LoginResponse{Jwt: loginToken}, nil
}

package user

import (
	"encoding/json"
	"fmt"
	jsonRpcClient "github.com/iot-my-world/brain/pkg/communication/jsonRpc/client"
	basicJsonRpcClient "github.com/iot-my-world/brain/pkg/communication/jsonRpc/client/basic"
	partyRegistrar "github.com/iot-my-world/brain/pkg/party/registrar"
	partyJsonRpcRegistrar "github.com/iot-my-world/brain/pkg/party/registrar/jsonRpc"
	"github.com/iot-my-world/brain/pkg/search/identifier/id"
	humanUserAdministrator "github.com/iot-my-world/brain/pkg/user/human/administrator"
	humanUserJsonRpcAdministrator "github.com/iot-my-world/brain/pkg/user/human/administrator/jsonRpc"
	authJsonRpcAdaptor "github.com/iot-my-world/brain/security/authorization/service/adaptor/jsonRpc"
	"github.com/iot-my-world/brain/security/claims"
	"github.com/iot-my-world/brain/security/claims/registerCompanyUser"
	wrappedClaims "github.com/iot-my-world/brain/security/claims/wrapped"
	companyTestData "github.com/iot-my-world/brain/testing/company/data"
	testData "github.com/iot-my-world/brain/testing/data"
	"github.com/stretchr/testify/suite"
	"gopkg.in/square/go-jose.v2"
	"reflect"
	"strings"
)

type User struct {
	suite.Suite
	jsonRpcClient          jsonRpcClient.Client
	humanUserAdministrator humanUserAdministrator.Administrator
	partyRegistrar         partyRegistrar.Registrar
}

func (suite *User) SetupTest() {
	// create the client
	suite.jsonRpcClient = basicJsonRpcClient.New(testData.BrainURL)

	suite.humanUserAdministrator = humanUserJsonRpcAdministrator.New(suite.jsonRpcClient)
	suite.partyRegistrar = partyJsonRpcRegistrar.New(suite.jsonRpcClient)
}

func (suite *User) TestCompanyInviteAndRegisterUsers() {
	// for each test company entity
	for companyIdx := range companyTestData.EntitiesAndAdminUsersToCreate {
		companyTestDataEntity := &companyTestData.EntitiesAndAdminUsersToCreate[companyIdx]
		// log in
		if err := suite.jsonRpcClient.Login(authJsonRpcAdaptor.LoginRequest{
			UsernameOrEmailAddress: companyTestDataEntity.AdminUser.Username,
			Password:               string(companyTestDataEntity.AdminUser.Password),
		}); err != nil {
			suite.FailNow(fmt.Sprintf("failed to log in as %s", companyTestDataEntity.AdminUser.Username), err.Error())
			return
		}

		// for each user assigned to this company
		for userIdx := range (*companyTestDataEntity).Users {
			// the minimal user must have an email address
			userEntity := &(*companyTestDataEntity).Users[userIdx]

			// the user has the same party details as the company admin user performing this invite

			// create minimal company user
			(*userEntity).ParentPartyType = suite.jsonRpcClient.Claims().PartyDetails().ParentPartyType
			(*userEntity).ParentId = suite.jsonRpcClient.Claims().PartyDetails().ParentId
			(*userEntity).PartyType = suite.jsonRpcClient.Claims().PartyDetails().PartyType
			(*userEntity).PartyId = suite.jsonRpcClient.Claims().PartyDetails().PartyId

			createCompanyUserResponse, err := suite.humanUserAdministrator.Create(&humanUserAdministrator.CreateRequest{
				User: *userEntity,
			})
			if err != nil {
				suite.FailNow("create company user failed", err.Error())
				return
			}
			// update id from created user
			(*userEntity).Id = createCompanyUserResponse.User.Id

			// invite the user
			inviteCompanyUserResponse, err := suite.partyRegistrar.InviteUser(&partyRegistrar.InviteUserRequest{
				UserIdentifier: id.Identifier{Id: (*userEntity).Id},
			})
			if err != nil {
				suite.FailNow("invite company user failed", err.Error())
				return
			}

			// parse the urlToken into a jsonWebToken object
			jwt := inviteCompanyUserResponse.URLToken[strings.Index(inviteCompanyUserResponse.URLToken, "&t=")+3:]
			jwtObject, err := jose.ParseSigned(jwt)
			if err != nil {
				suite.FailNow("error parsing jwt", err.Error())
			}

			// Access Underlying jwt payload bytes without verification
			jwtPayload := reflect.ValueOf(jwtObject).Elem().FieldByName("payload")

			// parse the bytes into wrapped claims
			wrapped := wrappedClaims.Wrapped{}
			if err := json.Unmarshal(jwtPayload.Bytes(), &wrapped); err != nil {
				suite.FailNow("error unmarshalling claims", err.Error())
			}

			// unwrap the claims into a claims.Claims interface
			unwrappedClaims, err := wrapped.Unwrap()
			if err != nil {
				suite.FailNow("error unwrapping claims", err.Error())
			}

			// confirm that the claims Type is correct
			if !suite.Equal(claims.RegisterCompanyUser, unwrappedClaims.Type(), "claims should be "+claims.RegisterCompanyUser) {
				suite.FailNow(fmt.Sprintf("claims are not of type %s", claims.RegisterCompanyUser))
			}

			// infer the interfaces type and update the client admin user entity with details from them
			switch typedClaims := unwrappedClaims.(type) {
			case registerCompanyUser.RegisterCompanyUser:
				(*userEntity).Id = typedClaims.User.Id
				(*userEntity).EmailAddress = typedClaims.User.EmailAddress
				(*userEntity).ParentPartyType = typedClaims.User.ParentPartyType
				(*userEntity).ParentId = typedClaims.User.ParentId
				(*userEntity).PartyType = typedClaims.User.PartyType
				(*userEntity).PartyId = typedClaims.User.PartyId
				// other userEntity fields already set in data for this test. Would have been filled out by user
			default:
				suite.FailNow(fmt.Sprintf("claims could not be inferred to type %s", claims.RegisterCompanyUser))
			}

			// store login token
			logInToken := suite.jsonRpcClient.GetJWT()
			// change token to registration token
			if err := suite.jsonRpcClient.SetJWT(jwt); err != nil {
				suite.FailNow("failed to set json rpc client jwt for registration", err.Error())
			}

			// register the company user
			// register the company admin user
			registerCompanyAdminUserResponse, err := suite.partyRegistrar.RegisterCompanyUser(&partyRegistrar.RegisterCompanyUserRequest{
				User: *userEntity,
			})
			if err != nil {
				suite.FailNow("error registering company user", err.Error())
				return
			}

			// set token back to logInToken
			if err := suite.jsonRpcClient.SetJWT(logInToken); err != nil {
				suite.FailNow("failed to set json rpc client jwt back to logInToken", err.Error())
			}

			// update the user with the response
			(*userEntity).Roles = registerCompanyAdminUserResponse.User.Roles
		}

		suite.jsonRpcClient.Logout()
	}
}

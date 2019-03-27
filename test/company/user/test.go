package user

import (
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/suite"
	jsonRpcClient "gitlab.com/iotTracker/brain/communication/jsonRpc/client"
	basicJsonRpcClient "gitlab.com/iotTracker/brain/communication/jsonRpc/client/basic"
	partyRegistrarJsonRpcAdaptor "gitlab.com/iotTracker/brain/party/registrar/adaptor/jsonRpc"
	"gitlab.com/iotTracker/brain/search/identifier/id"
	wrappedIdentifier "gitlab.com/iotTracker/brain/search/identifier/wrapped"
	authJsonRpcAdaptor "gitlab.com/iotTracker/brain/security/auth/service/adaptor/jsonRpc"
	"gitlab.com/iotTracker/brain/security/claims"
	"gitlab.com/iotTracker/brain/security/claims/registerCompanyUser"
	wrappedClaims "gitlab.com/iotTracker/brain/security/claims/wrapped"
	companyTestData "gitlab.com/iotTracker/brain/test/company/data"
	testData "gitlab.com/iotTracker/brain/test/data"
	userAdministratorJsonRpcAdaptor "gitlab.com/iotTracker/brain/user/administrator/adaptor/jsonRpc"
	"gopkg.in/square/go-jose.v2"
	"reflect"
	"strings"
)

type User struct {
	suite.Suite
	jsonRpcClient jsonRpcClient.Client
}

func (suite *User) SetupTest() {
	// create the client
	suite.jsonRpcClient = basicJsonRpcClient.New(testData.BrainURL)
}

func (suite *User) TestInviteAndRegisterUsers() {
	// for each test company entity
	for companyIdx := range companyTestData.EntitiesAndAdminUsersToCreate {
		companyTestDataEntity := &companyTestData.EntitiesAndAdminUsersToCreate[companyIdx]
		// log in
		if err := suite.jsonRpcClient.Login(authJsonRpcAdaptor.LoginRequest{
			UsernameOrEmailAddress: companyTestDataEntity.AdminUser.Username,
			Password:               string(companyTestDataEntity.AdminUser.Password),
		}); err != nil {
			suite.FailNow(fmt.Sprintf("failed to log in as %s", companyTestDataEntity.AdminUser.Username), err.Error())
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

			createCompanyUserResponse := userAdministratorJsonRpcAdaptor.CreateResponse{}
			if err := suite.jsonRpcClient.JsonRpcRequest(
				"UserAdministrator.Create",
				userAdministratorJsonRpcAdaptor.CreateRequest{
					User: *userEntity,
				},
				&createCompanyUserResponse,
			); err != nil {
				suite.FailNow("create company user failed", err.Error())
			}
			// update id
			(*userEntity).Id = createCompanyUserResponse.User.Id

			// create identifier for the user entity to invite
			userIdentifier, err := wrappedIdentifier.Wrap(id.Identifier{Id: (*userEntity).Id})
			if err != nil {
				suite.FailNow("error wrapping userIdentifier", err.Error())
			}

			// invite the user
			inviteCompanyUserResponse := partyRegistrarJsonRpcAdaptor.InviteUserResponse{}
			if err := suite.jsonRpcClient.JsonRpcRequest(
				"PartyRegistrar.InviteUser",
				partyRegistrarJsonRpcAdaptor.InviteUserRequest{
					UserIdentifier: *userIdentifier,
				},
				&inviteCompanyUserResponse,
			); err != nil {
				suite.FailNow("invite company user failed", err.Error())
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

			// create a new json rpc client to register the user with
			registerJsonRpcClient := basicJsonRpcClient.New(testData.BrainURL)
			if err := registerJsonRpcClient.SetJWT(jwt); err != nil {
				suite.FailNow("failed to set jwt in registration client", err.Error())
			}

			// register the company user
			registerCompanyResponse := partyRegistrarJsonRpcAdaptor.RegisterCompanyUserResponse{}
			if err := registerJsonRpcClient.JsonRpcRequest(
				"PartyRegistrar.RegisterCompanyUser",
				partyRegistrarJsonRpcAdaptor.RegisterCompanyUserRequest{
					User: *userEntity,
				},
				&registerCompanyResponse,
			); err != nil {
				suite.FailNow("error registering company user", err.Error())
			}

			// update the user with the response
			(*userEntity).Roles = registerCompanyResponse.User.Roles
		}

		suite.jsonRpcClient.Logout()
	}
}

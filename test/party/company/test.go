package company

import (
	"github.com/stretchr/testify/suite"
	jsonRpcClient "gitlab.com/iotTracker/brain/communication/jsonRpc/client"
	basicJsonRpcClient "gitlab.com/iotTracker/brain/communication/jsonRpc/client/basic"
	authJsonRpcAdaptor "gitlab.com/iotTracker/brain/security/auth/service/adaptor/jsonRpc"
	clientTest "gitlab.com/iotTracker/brain/test/party/client"
	clientRecordHandlerJsonRpcAdaptor "gitlab.com/iotTracker/brain/party/client/recordHandler/adaptor/jsonRpc"
	partyRegistrarJsonRpcAdaptor "gitlab.com/iotTracker/brain/party/registrar/adaptor/jsonRpc"
	"fmt"
	"gitlab.com/iotTracker/brain/search/identifier/id"
	"strings"
	"gopkg.in/square/go-jose.v2"
	"reflect"
	"gitlab.com/iotTracker/brain/security/wrappedClaims"
	"gitlab.com/iotTracker/brain/security/claims"
	"encoding/json"
	"gitlab.com/iotTracker/brain/security/claims/registerClientAdminUser"
	"gitlab.com/iotTracker/brain/party/user"
	"gitlab.com/iotTracker/brain/party"
)

type Company struct {
	suite.Suite
	jsonRpcClient jsonRpcClient.Client
}

func (suite *Company) SetupTest() {
	// create the client
	suite.jsonRpcClient = basicJsonRpcClient.New("http://localhost:9010/api")
}

func (suite *Company) TestCompanyCreateClients() {
	// for each test company entity
	for _, companyTestDataEntity := range EntitiesAndAdminUsersToCreate {
		// log in
		if err := suite.jsonRpcClient.Login(authJsonRpcAdaptor.LoginRequest{
			UsernameOrEmailAddress: companyTestDataEntity.AdminUser.Username,
			Password:               string(companyTestDataEntity.AdminUser.Password),
		}); err != nil {
			suite.FailNow(fmt.Sprintf("failed to log in as %s", companyTestDataEntity.AdminUser.Username), err.Error())
		}

		// for each client assigned to be owned by this company
		for idx := range clientTest.EntitiesAndAdminUsersToCreate[companyTestDataEntity.Company.Name] {
			clientEntity := &clientTest.EntitiesAndAdminUsersToCreate[companyTestDataEntity.Company.Name][idx].Client

			// update the entity
			(*clientEntity).ParentPartyType = suite.jsonRpcClient.Claims().PartyDetails().PartyType
			(*clientEntity).ParentId = suite.jsonRpcClient.Claims().PartyDetails().PartyId

			// create the client
			clientCreateResponse := clientRecordHandlerJsonRpcAdaptor.CreateResponse{}
			if err := suite.jsonRpcClient.JsonRpcRequest(
				"ClientRecordHandler.Create",
				clientRecordHandlerJsonRpcAdaptor.CreateRequest{
					Client: *clientEntity,
				},
				&clientCreateResponse,
			); err != nil {
				suite.FailNow("create client failed", err.Error())
			}

			// update the client
			(*clientEntity).Id = clientCreateResponse.Client.Id
		}

		// log out
		suite.jsonRpcClient.Logout()
	}
}

func (suite *Company) TestCompanyInviteAndRegisterClients() {
	// for each test company entity
	for _, companyTestDataEntity := range EntitiesAndAdminUsersToCreate {
		// log in
		if err := suite.jsonRpcClient.Login(authJsonRpcAdaptor.LoginRequest{
			UsernameOrEmailAddress: companyTestDataEntity.AdminUser.Username,
			Password:               string(companyTestDataEntity.AdminUser.Password),
		}); err != nil {
			suite.FailNow(fmt.Sprintf("failed to log in as %s", companyTestDataEntity.AdminUser.Username), err.Error())
		}

		// for each client assigned to be owned by this company
		for idx := range clientTest.EntitiesAndAdminUsersToCreate[companyTestDataEntity.Company.Name] {
			clientEntity := &clientTest.EntitiesAndAdminUsersToCreate[companyTestDataEntity.Company.Name][idx].Client

			// make minimal client admin user
			clientAdminUser := user.User{
				EmailAddress:    clientEntity.AdminEmailAddress,
				ParentPartyType: suite.jsonRpcClient.Claims().PartyDetails().PartyType,
				ParentId:        suite.jsonRpcClient.Claims().PartyDetails().ParentId,
				PartyType:       party.Client,
				PartyId:         id.Identifier{Id: (*clientEntity).Id},
			}

			// invite the admin user
			inviteClientAdminUserResponse := partyRegistrarJsonRpcAdaptor.InviteClientAdminUserResponse{}
			if err := suite.jsonRpcClient.JsonRpcRequest(
				"PartyRegistrar.InviteClientAdminUser",
				partyRegistrarJsonRpcAdaptor.InviteCompanyAdminUserRequest{
					User: clientAdminUser,
				},
				&inviteClientAdminUserResponse,
			); err != nil {
				suite.FailNow("invite client admin user failed", err.Error())
			}

			// parse the urlToken into a jsonWebToken object
			jwt := inviteClientAdminUserResponse.URLToken[strings.Index(inviteClientAdminUserResponse.URLToken, "&t=")+3:]
			jwtObject, err := jose.ParseSigned(jwt)
			if err != nil {
				suite.FailNow("error parsing jwt", err.Error())
			}

			// Access Underlying jwt payload bytes without verification
			jwtPayload := reflect.ValueOf(jwtObject).Elem().FieldByName("payload")

			// parse the bytes into wrapped claims
			wrapped := wrappedClaims.WrappedClaims{}
			if err := json.Unmarshal(jwtPayload.Bytes(), &wrapped); err != nil {
				suite.FailNow("error unmarshalling claims", err.Error())
			}

			// unwrap the claims into a claims.Claims interface
			unwrappedClaims, err := wrapped.Unwrap()
			if err != nil {
				suite.FailNow("error unwrapping claims", err.Error())
			}

			// confirm that the claims Type is correct
			if !suite.Equal(claims.RegisterClientAdminUser, unwrappedClaims.Type(), "claims should be "+claims.RegisterClientAdminUser) {
				suite.FailNow(fmt.Sprintf("claims are not of type %s", claims.RegisterClientAdminUser))
			}

			// infer the interfaces type and update the client admin user entity with details from them
			clientAdminUserEntity := &clientTest.EntitiesAndAdminUsersToCreate[companyTestDataEntity.Company.Name][idx].AdminUser
			switch typedClaims := unwrappedClaims.(type) {
			case registerClientAdminUser.RegisterClientAdminUser:
				(*clientAdminUserEntity).EmailAddress = typedClaims.EmailAddress
				(*clientAdminUserEntity).ParentPartyType = typedClaims.ParentPartyType
				(*clientAdminUserEntity).ParentId = typedClaims.ParentId
				(*clientAdminUserEntity).PartyType = typedClaims.PartyType
				(*clientAdminUserEntity).PartyId = typedClaims.PartyId
			default:
				suite.FailNow(fmt.Sprintf("claims could not be inferred to type %s", claims.RegisterClientAdminUser))
			}

			// create a new json rpc client to register the user with
			registerJsonRpcClient := basicJsonRpcClient.New("http://localhost:9010/api")
			if err := registerJsonRpcClient.SetJWT(jwt); err != nil {
				suite.FailNow("failed to set jwt in registration client", err.Error())
			}

			// register the client admin user
			registerClientAdminUserResponse := partyRegistrarJsonRpcAdaptor.RegisterClientAdminUserResponse{}
			password := string(clientAdminUserEntity.Password)
			(*clientAdminUserEntity).Password = []byte{}
			if err := registerJsonRpcClient.JsonRpcRequest(
				"PartyRegistrar.RegisterClientAdminUser",
				partyRegistrarJsonRpcAdaptor.RegisterClientAdminUserRequest{
					User:     *clientAdminUserEntity,
					Password: password,
				},
				&registerClientAdminUserResponse,
			); err != nil {
				suite.FailNow("error registering client admin user", err.Error())
			}

			// update the client admin user entity
			(*clientAdminUserEntity).Id = registerClientAdminUserResponse.User.Id
			(*clientAdminUserEntity).Roles = registerClientAdminUserResponse.User.Roles
			(*clientAdminUserEntity).Password = []byte(password)
		}

		// log out
		suite.jsonRpcClient.Logout()
	}
}

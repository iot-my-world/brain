package client

import (
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/suite"
	jsonRpcClient "gitlab.com/iotTracker/brain/communication/jsonRpc/client"
	basicJsonRpcClient "gitlab.com/iotTracker/brain/communication/jsonRpc/client/basic"
	clientAdministratorJsonRpcAdaptor "gitlab.com/iotTracker/brain/party/client/administrator/adaptor/jsonRpc"
	partyRegistrarJsonRpcAdaptor "gitlab.com/iotTracker/brain/party/registrar/adaptor/jsonRpc"
	"gitlab.com/iotTracker/brain/search/identifier/id"
	wrappedIdentifier "gitlab.com/iotTracker/brain/search/identifier/wrapped"
	authJsonRpcAdaptor "gitlab.com/iotTracker/brain/security/auth/service/adaptor/jsonRpc"
	"gitlab.com/iotTracker/brain/security/claims"
	"gitlab.com/iotTracker/brain/security/claims/registerClientAdminUser"
	wrappedClaims "gitlab.com/iotTracker/brain/security/claims/wrapped"
	clientTestData "gitlab.com/iotTracker/brain/test/client/data"
	companyTestData "gitlab.com/iotTracker/brain/test/company/data"
	testData "gitlab.com/iotTracker/brain/test/data"
	"gopkg.in/square/go-jose.v2"
	"reflect"
	"strings"
)

type Client struct {
	suite.Suite
	jsonRpcClient jsonRpcClient.Client
}

func (suite *Client) SetupTest() {
	// create the client
	suite.jsonRpcClient = basicJsonRpcClient.New(testData.BrainURL)
}

func (suite *Client) TestCompanyCreateClients() {
	// for each test company entity
	for _, companyTestDataEntity := range companyTestData.EntitiesAndAdminUsersToCreate {
		// log in
		if err := suite.jsonRpcClient.Login(authJsonRpcAdaptor.LoginRequest{
			UsernameOrEmailAddress: companyTestDataEntity.AdminUser.Username,
			Password:               string(companyTestDataEntity.AdminUser.Password),
		}); err != nil {
			suite.FailNow(fmt.Sprintf("failed to log in as %s", companyTestDataEntity.AdminUser.Username), err.Error())
		}

		// for each client assigned to be owned by this company
		for idx := range clientTestData.EntitiesAndAdminUsersToCreate[companyTestDataEntity.Company.Name] {
			clientEntity := &clientTestData.EntitiesAndAdminUsersToCreate[companyTestDataEntity.Company.Name][idx].Client

			// update the entity
			(*clientEntity).ParentPartyType = suite.jsonRpcClient.Claims().PartyDetails().PartyType
			(*clientEntity).ParentId = suite.jsonRpcClient.Claims().PartyDetails().PartyId

			// create the client
			clientCreateResponse := clientAdministratorJsonRpcAdaptor.CreateResponse{}
			if err := suite.jsonRpcClient.JsonRpcRequest(
				"ClientAdministrator.Create",
				clientAdministratorJsonRpcAdaptor.CreateRequest{
					Client: *clientEntity,
				},
				&clientCreateResponse,
			); err != nil {
				suite.FailNow("create client failed", err.Error())
			}

			// update the client
			clientEntity.Id = clientCreateResponse.Client.Id
		}

		// log out
		suite.jsonRpcClient.Logout()
	}
}

func (suite *Client) TestCompanyInviteAndRegisterClients() {
	// for each test company entity
	for _, companyTestDataEntity := range companyTestData.EntitiesAndAdminUsersToCreate {
		// log in
		if err := suite.jsonRpcClient.Login(authJsonRpcAdaptor.LoginRequest{
			UsernameOrEmailAddress: companyTestDataEntity.AdminUser.Username,
			Password:               string(companyTestDataEntity.AdminUser.Password),
		}); err != nil {
			suite.FailNow(fmt.Sprintf("failed to log in as %s", companyTestDataEntity.AdminUser.Username), err.Error())
		}

		// for each client assigned to be owned by this company
		for idx := range clientTestData.EntitiesAndAdminUsersToCreate[companyTestDataEntity.Company.Name] {
			clientEntity := &clientTestData.EntitiesAndAdminUsersToCreate[companyTestDataEntity.Company.Name][idx].Client
			clientAdminUserEntity := &clientTestData.EntitiesAndAdminUsersToCreate[companyTestDataEntity.Company.Name][idx].AdminUser

			// create identifier for the client entity
			clientIdentifier, err := wrappedIdentifier.Wrap(id.Identifier{Id: clientEntity.Id})
			if err != nil {
				suite.FailNow("error wrapping client identifier", err.Error())
			}

			// invite the admin user
			inviteClientAdminUserResponse := partyRegistrarJsonRpcAdaptor.InviteClientAdminUserResponse{}
			if err := suite.jsonRpcClient.JsonRpcRequest(
				"PartyRegistrar.InviteClientAdminUser",
				partyRegistrarJsonRpcAdaptor.InviteClientAdminUserRequest{
					WrappedClientIdentifier: *clientIdentifier,
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
			if !suite.Equal(claims.RegisterClientAdminUser, unwrappedClaims.Type(), "claims should be "+claims.RegisterClientAdminUser) {
				suite.FailNow(fmt.Sprintf("claims are not of type %s", claims.RegisterClientAdminUser))
			}

			// infer the interfaces type and update the client admin user entity with details from them
			switch typedClaims := unwrappedClaims.(type) {
			case registerClientAdminUser.RegisterClientAdminUser:
				(*clientAdminUserEntity).Id = typedClaims.User.Id
				(*clientAdminUserEntity).EmailAddress = typedClaims.User.EmailAddress
				(*clientAdminUserEntity).ParentPartyType = typedClaims.User.ParentPartyType
				(*clientAdminUserEntity).ParentId = typedClaims.User.ParentId
				(*clientAdminUserEntity).PartyType = typedClaims.User.PartyType
				(*clientAdminUserEntity).PartyId = typedClaims.User.PartyId
			default:
				suite.FailNow(fmt.Sprintf("claims could not be inferred to type %s", claims.RegisterClientAdminUser))
			}

			// create a new json rpc client to register the user with
			registerJsonRpcClient := basicJsonRpcClient.New(testData.BrainURL)
			if err := registerJsonRpcClient.SetJWT(jwt); err != nil {
				suite.FailNow("failed to set jwt in registration client", err.Error())
			}

			// register the client admin user
			registerClientAdminUserResponse := partyRegistrarJsonRpcAdaptor.RegisterClientAdminUserResponse{}
			if err := registerJsonRpcClient.JsonRpcRequest(
				"PartyRegistrar.RegisterClientAdminUser",
				partyRegistrarJsonRpcAdaptor.RegisterClientAdminUserRequest{
					User: *clientAdminUserEntity,
				},
				&registerClientAdminUserResponse,
			); err != nil {
				suite.FailNow("error registering client admin user", err.Error())
			}

			// update the client admin user entity
			(*clientAdminUserEntity).Id = registerClientAdminUserResponse.User.Id
			(*clientAdminUserEntity).Roles = registerClientAdminUserResponse.User.Roles
		}

		// log out
		suite.jsonRpcClient.Logout()
	}
}

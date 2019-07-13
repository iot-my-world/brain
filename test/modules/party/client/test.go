package client

import (
	"encoding/json"
	"fmt"
	jsonRpcClient "github.com/iot-my-world/brain/pkg/api/jsonRpc/client"
	basicJsonRpcClient "github.com/iot-my-world/brain/pkg/api/jsonRpc/client/basic"
	jsonRpcServerAuthenticator "github.com/iot-my-world/brain/pkg/api/jsonRpc/server/authenticator"
	"github.com/iot-my-world/brain/pkg/party/client"
	clientAdministrator "github.com/iot-my-world/brain/pkg/party/client/administrator"
	clientJsonRpcAdministrator "github.com/iot-my-world/brain/pkg/party/client/administrator/jsonRpc"
	clientRecordHandler "github.com/iot-my-world/brain/pkg/party/client/recordHandler"
	clientJsonRpcRecordHandler "github.com/iot-my-world/brain/pkg/party/client/recordHandler/jsonRpc"
	partyRegistrar "github.com/iot-my-world/brain/pkg/party/registrar"
	partyJsonRpcRegistrar "github.com/iot-my-world/brain/pkg/party/registrar/jsonRpc"
	"github.com/iot-my-world/brain/pkg/search/criterion"
	"github.com/iot-my-world/brain/pkg/search/identifier/adminEmailAddress"
	"github.com/iot-my-world/brain/pkg/search/identifier/emailAddress"
	"github.com/iot-my-world/brain/pkg/search/identifier/id"
	"github.com/iot-my-world/brain/pkg/search/query"
	"github.com/iot-my-world/brain/pkg/security/claims"
	"github.com/iot-my-world/brain/pkg/security/claims/registerClientAdminUser"
	"github.com/iot-my-world/brain/pkg/security/claims/registerClientUser"
	wrappedClaims "github.com/iot-my-world/brain/pkg/security/claims/wrapped"
	humanUser "github.com/iot-my-world/brain/pkg/user/human"
	humanUserAdministrator "github.com/iot-my-world/brain/pkg/user/human/administrator"
	humanUserJsonRpcAdministrator "github.com/iot-my-world/brain/pkg/user/human/administrator/jsonRpc"
	humanUserRecordHandler "github.com/iot-my-world/brain/pkg/user/human/recordHandler"
	humanUserJsonRpcRecordHandler "github.com/iot-my-world/brain/pkg/user/human/recordHandler/jsonRpc"
	"github.com/stretchr/testify/suite"
	"gopkg.in/square/go-jose.v2"
	"reflect"
	"strings"
)

func New(
	url string,
	user humanUser.User,
	testData []Data,
) *test {
	return &test{
		testData:      testData,
		user:          user,
		jsonRpcClient: basicJsonRpcClient.New(url),
	}
}

type test struct {
	suite.Suite
	jsonRpcClient          jsonRpcClient.Client
	clientRecordHandler    clientRecordHandler.RecordHandler
	clientAdministrator    clientAdministrator.Administrator
	humanUserAdministrator humanUserAdministrator.Administrator
	humanUserRecordHandler humanUserRecordHandler.RecordHandler
	partyRegistrar         partyRegistrar.Registrar
	user                   humanUser.User
	testData               []Data
}

type Data struct {
	Client    client.Client
	AdminUser humanUser.User
	Users     []humanUser.User
}

func (suite *test) SetupTest() {

	// log in the client
	if err := suite.jsonRpcClient.Login(jsonRpcServerAuthenticator.LoginRequest{
		UsernameOrEmailAddress: suite.user.Username,
		Password:               string(suite.user.Password),
	}); err != nil {
		suite.Fail("log in error", err.Error())
		return
	}

	// set up service provider clients that use jsonRpcClient
	suite.clientRecordHandler = clientJsonRpcRecordHandler.New(suite.jsonRpcClient)
	suite.clientAdministrator = clientJsonRpcAdministrator.New(suite.jsonRpcClient)
	suite.partyRegistrar = partyJsonRpcRegistrar.New(suite.jsonRpcClient)
	suite.humanUserAdministrator = humanUserJsonRpcAdministrator.New(suite.jsonRpcClient)
	suite.humanUserRecordHandler = humanUserJsonRpcRecordHandler.New(suite.jsonRpcClient)
}

func (suite *test) TestClient1Create() {
	// create all clients in test data
	for _, data := range suite.testData {
		clientEntity := data.Client

		// update the new client's details as would be done from the front end
		clientEntity.ParentPartyType = suite.jsonRpcClient.Claims().PartyDetails().PartyType
		clientEntity.ParentId = suite.jsonRpcClient.Claims().PartyDetails().PartyId

		// create the client
		if _, err := suite.clientAdministrator.Create(&clientAdministrator.CreateRequest{
			Client: clientEntity,
		}); err != nil {
			suite.FailNow("create client failed", err.Error())
			return
		}
	}

	// collect all clients
	clientCollectResponse, err := suite.clientRecordHandler.Collect(&clientRecordHandler.CollectRequest{
		Criteria: make([]criterion.Criterion, 0),
		Query:    query.Query{},
	})
	if err != nil {
		suite.Failf("collect clients failed", err.Error())
		return
	}

	// confirm that each created client can be found
nextClientToCreate:
	// for every client that should be created
	for _, clientToCreate := range suite.testData {
		// look for clientToCreate among collected clients
		for _, existingClient := range clientCollectResponse.Records {
			if clientToCreate.Client.AdminEmailAddress == existingClient.AdminEmailAddress {
				// update fields set during creation
				clientToCreate.Client.Id = existingClient.Id
				clientToCreate.Client.ParentPartyType = existingClient.ParentPartyType
				clientToCreate.Client.ParentId = existingClient.ParentId

				suite.Equal(
					clientToCreate.Client,
					existingClient,
					"created client should be equal",
				)
				// if it is found and equal, check for next client to create
				continue nextClientToCreate
			}
		}
		// if execution reaches here then clientToCreate was not found among collected clients
	}
}

func (suite *test) TestClient2UpdateAllowedFields() {
	for _, data := range suite.testData {

		// retrieve the client by admin email address
		clientRetrieveResponse, err := suite.clientRecordHandler.Retrieve(&clientRecordHandler.RetrieveRequest{
			Identifier: adminEmailAddress.Identifier{
				AdminEmailAddress: data.Client.AdminEmailAddress,
			},
		})
		if err != nil {
			suite.FailNow("retrieve client entity failed", err.Error())
			return
		}

		// copy retrieved client
		updatedClientEntity := clientRetrieveResponse.Client

		// update allowed fields
		updatedClientEntity.Name = "Changed Name"

		// perform update
		updateAllowedFieldsResponse, err := suite.clientAdministrator.UpdateAllowedFields(&clientAdministrator.UpdateAllowedFieldsRequest{
			Client: updatedClientEntity,
		})
		if err != nil {
			suite.FailNow("client update allowed fields failed", err.Error())
			return
		}

		suite.Equal(
			updatedClientEntity,
			updateAllowedFieldsResponse.Client,
			"updated client should equal client in updated response",
		)

		// retrieve the updated entity by id
		updatedClientRetrieveResponse, err := suite.clientRecordHandler.Retrieve(&clientRecordHandler.RetrieveRequest{
			Identifier: id.Identifier{
				Id: updatedClientEntity.Id,
			},
		})
		if err != nil {
			suite.FailNow("retrieve updated client entity failed", err.Error())
			return
		}

		suite.Equal(
			updatedClientEntity,
			updatedClientRetrieveResponse.Client,
			"retrieved client should equal updated client",
		)

		// update client back to original
		updateAllowedFieldsResponse, err = suite.clientAdministrator.UpdateAllowedFields(&clientAdministrator.UpdateAllowedFieldsRequest{
			Client: clientRetrieveResponse.Client,
		})
		if err != nil {
			suite.FailNow("client update allowed fields failed", err.Error())
			return
		}

		suite.Equal(
			clientRetrieveResponse.Client,
			updateAllowedFieldsResponse.Client,
			"updated client should equal client in updated response",
		)
	}
}

func (suite *test) TestClient3Delete() {
	// create a client
	createResponse, err := suite.clientAdministrator.Create(&clientAdministrator.CreateRequest{
		Client: client.Client{
			Type:              client.Individual,
			Name:              "BobToBeDeleted",
			AdminEmailAddress: "bob@gmail.com",
			ParentPartyType:   suite.jsonRpcClient.Claims().PartyDetails().PartyType,
			ParentId:          suite.jsonRpcClient.Claims().PartyDetails().PartyId,
		},
	})
	if err != nil {
		suite.FailNow("error creating client", err.Error())
		return
	}

	// retrieve the client
	retrieveResponse, err := suite.clientRecordHandler.Retrieve(&clientRecordHandler.RetrieveRequest{
		Identifier: id.Identifier{
			Id: createResponse.Client.Id,
		},
	})
	if err != nil {
		suite.FailNow("error retrieving client", err.Error())
		return
	}

	// delete the client
	if _, err := suite.clientAdministrator.Delete(&clientAdministrator.DeleteRequest{
		ClientIdentifier: id.Identifier{
			Id: retrieveResponse.Client.Id,
		},
	}); err != nil {
		suite.FailNow("error deleting client", err.Error())
		return
	}

	// collect all clients
	collectResponse, err := suite.clientRecordHandler.Collect(&clientRecordHandler.CollectRequest{
		Criteria: make([]criterion.Criterion, 0),
		Query:    query.Query{},
	})
	if err != nil {
		suite.FailNow("error collecting clients", err.Error())
		return
	}

	// confirm that deleted client not among collected clients
	for _, c := range collectResponse.Records {
		if c.Id == retrieveResponse.Client.Id {
			suite.FailNow("client found in collected clients after deletion")
		}
	}
}

func (suite *test) TestClient4InviteAndRegisterAdmin() {
	for _, data := range suite.testData {
		clientEntity := data.Client
		clientAdminUserEntity := data.AdminUser

		// retrieve the client by admin email address
		clientRetrieveResponse, err := suite.clientRecordHandler.Retrieve(&clientRecordHandler.RetrieveRequest{
			Identifier: adminEmailAddress.Identifier{
				AdminEmailAddress: clientEntity.AdminEmailAddress,
			},
		})
		if err != nil {
			suite.FailNow("retrieve client entity failed", err.Error())
			return
		}

		// invite the admin user
		inviteClientAdminUserResponse, err := suite.partyRegistrar.InviteClientAdminUser(&partyRegistrar.InviteClientAdminUserRequest{
			ClientIdentifier: id.Identifier{Id: clientRetrieveResponse.Client.Id},
		})
		if err != nil {
			suite.FailNow("invite client admin user failed", err.Error())
			return
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
			return
		}

		// confirm that the claims Type is correct
		if !suite.Equal(claims.RegisterClientAdminUser, unwrappedClaims.Type(), "claims should be "+claims.RegisterClientAdminUser) {
			suite.FailNow(fmt.Sprintf("claims are not of type %s", claims.RegisterClientAdminUser))
		}

		// infer the interfaces type and update the client admin user entity with details from them
		switch typedClaims := unwrappedClaims.(type) {
		case registerClientAdminUser.RegisterClientAdminUser:
			clientAdminUserEntity.Id = typedClaims.User.Id
			clientAdminUserEntity.EmailAddress = typedClaims.User.EmailAddress
			clientAdminUserEntity.ParentPartyType = typedClaims.User.ParentPartyType
			clientAdminUserEntity.ParentId = typedClaims.User.ParentId
			clientAdminUserEntity.PartyType = typedClaims.User.PartyType
			clientAdminUserEntity.PartyId = typedClaims.User.PartyId
		default:
			suite.FailNow(fmt.Sprintf("claims could not be inferred to type %s", claims.RegisterClientAdminUser))
		}

		// store login token
		logInToken := suite.jsonRpcClient.GetJWT()
		// change token to registration token
		if err := suite.jsonRpcClient.SetJWT(jwt); err != nil {
			suite.FailNow("failed to set json rpc client jwt for registration", err.Error())
		}

		// register the client admin user
		if _, err := suite.partyRegistrar.RegisterClientAdminUser(&partyRegistrar.RegisterClientAdminUserRequest{
			User: clientAdminUserEntity,
		}); err != nil {
			suite.FailNow("error registering client admin user", err.Error())
			return
		}

		// set token back to logInToken
		if err := suite.jsonRpcClient.SetJWT(logInToken); err != nil {
			suite.FailNow("failed to set json rpc client jwt back to logInToken", err.Error())
		}
	}
}

func (suite *test) TestClient5CreateUsers() {
	for _, clientData := range suite.testData {
		// authenticate json rpc client as client admin user
		if err := suite.jsonRpcClient.Login(jsonRpcServerAuthenticator.LoginRequest{
			UsernameOrEmailAddress: clientData.AdminUser.Username,
			Password:               string(clientData.AdminUser.Password),
		}); err != nil {
			suite.FailNow("could not log in as client admin user", err.Error())
			return
		}

		for _, userToCreate := range clientData.Users {
			// set user's party details
			userToCreate.ParentPartyType = suite.jsonRpcClient.Claims().PartyDetails().ParentPartyType
			userToCreate.ParentId = suite.jsonRpcClient.Claims().PartyDetails().ParentId
			userToCreate.PartyType = suite.jsonRpcClient.Claims().PartyDetails().PartyType
			userToCreate.PartyId = suite.jsonRpcClient.Claims().PartyDetails().PartyId

			// create user
			createResponse, err := suite.humanUserAdministrator.Create(&humanUserAdministrator.CreateRequest{
				User: userToCreate,
			})
			if err != nil {
				suite.FailNow("error creating company user")
				return
			}
			// set fields set on creation
			userToCreate.Id = createResponse.User.Id
			if !suite.Equal(
				userToCreate,
				createResponse.User,
				"user in create response should be equal to user to create",
			) {
				return
			}

			// retrieve user
			retrieveUserResponse, err := suite.humanUserRecordHandler.Retrieve(&humanUserRecordHandler.RetrieveRequest{
				Identifier: emailAddress.Identifier{
					EmailAddress: userToCreate.EmailAddress,
				},
			})
			if err != nil {
				suite.FailNow("error retrieving user", err.Error())
				return
			}
			if !suite.Equal(
				userToCreate,
				retrieveUserResponse.User,
				"retrieved user should be the same as created",
			) {
				return
			}
		}
	}
}

func (suite *test) TestClient6InviteAndRegisterUsers() {
	for _, clientData := range suite.testData {
		// authenticate json rpc client as client admin user
		if err := suite.jsonRpcClient.Login(jsonRpcServerAuthenticator.LoginRequest{
			UsernameOrEmailAddress: clientData.AdminUser.Username,
			Password:               string(clientData.AdminUser.Password),
		}); err != nil {
			suite.FailNow("could not log in as client admin user", err.Error())
			return
		}

		for _, userToInvite := range clientData.Users {
			// invite user
			inviteUserResponse, err := suite.partyRegistrar.InviteUser(&partyRegistrar.InviteUserRequest{
				UserIdentifier: emailAddress.Identifier{
					EmailAddress: userToInvite.EmailAddress,
				},
			})
			if err != nil {
				suite.FailNow("invite client user failed", err.Error())
				return
			}
			// parse the urlToken into a jsonWebToken object
			jwt := inviteUserResponse.URLToken[strings.Index(inviteUserResponse.URLToken, "&t=")+3:]
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
				return
			}

			// confirm that the claims Type is correct
			if !suite.Equal(
				claims.RegisterClientUser,
				unwrappedClaims.Type(),
				"claims should be "+claims.RegisterClientUser,
			) {
				return
			}

			// infer the interfaces type and update the client user entity with details from them
			switch typedClaims := unwrappedClaims.(type) {
			case registerClientUser.RegisterClientUser:
				userToInvite.Id = typedClaims.User.Id
				userToInvite.EmailAddress = typedClaims.User.EmailAddress
				userToInvite.ParentPartyType = typedClaims.User.ParentPartyType
				userToInvite.ParentId = typedClaims.User.ParentId
				userToInvite.PartyType = typedClaims.User.PartyType
				userToInvite.PartyId = typedClaims.User.PartyId
			default:
				suite.FailNow(fmt.Sprintf("claims could not be inferred to type %s", claims.RegisterClientUser))
				return
			}

			// store login token
			logInToken := suite.jsonRpcClient.GetJWT()

			// change token to registration token
			if err := suite.jsonRpcClient.SetJWT(jwt); err != nil {
				suite.FailNow("failed to set json rpc client jwt for registration", err.Error())
			}

			// register the client user
			if _, err := suite.partyRegistrar.RegisterClientUser(&partyRegistrar.RegisterClientUserRequest{
				User: userToInvite,
			}); err != nil {
				suite.FailNow("error registering client user", err.Error())
				return
			}

			// set token back to logInToken
			if err := suite.jsonRpcClient.SetJWT(logInToken); err != nil {
				suite.FailNow("failed to set json rpc client jwt back to logInToken", err.Error())
			}
		}
	}
}

func (suite *test) TestClient7UserLogin() {
	for _, clientData := range suite.testData {
		for _, userToTest := range clientData.Users {
			if err := suite.jsonRpcClient.Login(jsonRpcServerAuthenticator.LoginRequest{
				UsernameOrEmailAddress: userToTest.Username,
				Password:               string(userToTest.Password),
			}); err != nil {
				suite.FailNow("could not log client user", err.Error())
				return
			}
		}
	}
}

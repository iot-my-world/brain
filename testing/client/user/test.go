package user

import (
	"encoding/json"
	"fmt"
	jsonRpcClient "github.com/iot-my-world/brain/communication/jsonRpc/client"
	basicJsonRpcClient "github.com/iot-my-world/brain/communication/jsonRpc/client/basic"
	partyRegistrarJsonRpcAdaptor "github.com/iot-my-world/brain/pkg/party/registrar/adaptor/jsonRpc"
	"github.com/iot-my-world/brain/search/identifier/id"
	wrappedIdentifier "github.com/iot-my-world/brain/search/identifier/wrapped"
	authJsonRpcAdaptor "github.com/iot-my-world/brain/security/authorization/service/adaptor/jsonRpc"
	"github.com/iot-my-world/brain/security/claims"
	"github.com/iot-my-world/brain/security/claims/registerClientUser"
	wrappedClaims "github.com/iot-my-world/brain/security/claims/wrapped"
	clientTestData "github.com/iot-my-world/brain/testing/client/data"
	testData "github.com/iot-my-world/brain/testing/data"
	userAdministratorJsonRpcAdaptor "github.com/iot-my-world/brain/user/human/administrator/adaptor/jsonRpc"
	"github.com/stretchr/testify/suite"
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

func (suite *User) TestClientInviteAndRegisterUsers() {
	for companyOwner := range clientTestData.EntitiesAndAdminUsersToCreate {
		for clientDataEntityIdx := range clientTestData.EntitiesAndAdminUsersToCreate[companyOwner] {
			clientDataEntity := &clientTestData.EntitiesAndAdminUsersToCreate[companyOwner][clientDataEntityIdx]
			// log in the client
			if err := suite.jsonRpcClient.Login(authJsonRpcAdaptor.LoginRequest{
				UsernameOrEmailAddress: clientDataEntity.AdminUser.Username,
				Password:               string(clientDataEntity.AdminUser.Password),
			}); err != nil {
				suite.FailNow(fmt.Sprintf("failed to log in as %s", clientDataEntity.AdminUser.Username), err.Error())
			}

			// invite and register all of the users
			for userIdx := range (*clientDataEntity).Users {
				userEntity := &(*clientDataEntity).Users[userIdx]

				// make minimal client user
				(*userEntity).ParentPartyType = suite.jsonRpcClient.Claims().PartyDetails().ParentPartyType
				(*userEntity).ParentId = suite.jsonRpcClient.Claims().PartyDetails().ParentId
				(*userEntity).PartyType = suite.jsonRpcClient.Claims().PartyDetails().PartyType
				(*userEntity).PartyId = suite.jsonRpcClient.Claims().PartyDetails().PartyId

				// create the user
				createCompanyUserResponse := userAdministratorJsonRpcAdaptor.CreateResponse{}
				if err := suite.jsonRpcClient.JsonRpcRequest(
					"UserAdministrator.Create",
					userAdministratorJsonRpcAdaptor.CreateRequest{
						User: *userEntity,
					},
					&createCompanyUserResponse,
				); err != nil {
					suite.FailNow("create client user failed", err.Error())
				}
				// update id
				(*userEntity).Id = createCompanyUserResponse.User.Id

				// create identifier for the user entity to invite
				userIdentifier, err := wrappedIdentifier.Wrap(id.Identifier{Id: (*userEntity).Id})
				if err != nil {
					suite.FailNow("error wrapping userIdentifier", err.Error())
				}

				// invite the user
				inviteClientUserResponse := partyRegistrarJsonRpcAdaptor.InviteUserResponse{}
				if err := suite.jsonRpcClient.JsonRpcRequest(
					"PartyRegistrar.InviteUser",
					partyRegistrarJsonRpcAdaptor.InviteUserRequest{
						WrappedUserIdentifier: *userIdentifier,
					},
					&inviteClientUserResponse,
				); err != nil {
					suite.FailNow("invite client user failed", err.Error())
				}

				// parse the urlToken into a jsonWebToken object
				jwt := inviteClientUserResponse.URLToken[strings.Index(inviteClientUserResponse.URLToken, "&t=")+3:]
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
				if !suite.Equal(claims.RegisterClientUser, unwrappedClaims.Type(), "claims should be "+claims.RegisterClientUser) {
					suite.FailNow(fmt.Sprintf("claims are not of type %s", claims.RegisterClientUser))
				}

				// infer the interfaces type and update the client admin user entity with details from them
				switch typedClaims := unwrappedClaims.(type) {
				case registerClientUser.RegisterClientUser:
					(*userEntity).Id = typedClaims.User.Id
					(*userEntity).EmailAddress = typedClaims.User.EmailAddress
					(*userEntity).ParentPartyType = typedClaims.User.ParentPartyType
					(*userEntity).ParentId = typedClaims.User.ParentId
					(*userEntity).PartyType = typedClaims.User.PartyType
					(*userEntity).PartyId = typedClaims.User.PartyId
				default:
					suite.FailNow(fmt.Sprintf("claims could not be inferred to type %s", claims.RegisterClientUser))
				}

				// create a new json rpc client to register the user with
				registerJsonRpcClient := basicJsonRpcClient.New(testData.BrainURL)
				if err := registerJsonRpcClient.SetJWT(jwt); err != nil {
					suite.FailNow("failed to set jwt in registration client", err.Error())
				}

				// register the client user
				registerClientResponse := partyRegistrarJsonRpcAdaptor.RegisterClientUserResponse{}
				if err := registerJsonRpcClient.JsonRpcRequest(
					"PartyRegistrar.RegisterClientUser",
					partyRegistrarJsonRpcAdaptor.RegisterClientUserRequest{
						User: *userEntity,
					},
					&registerClientResponse,
				); err != nil {
					suite.FailNow("error registering client user", err.Error())
				}

				// update the user
				(*userEntity).Roles = registerClientResponse.User.Roles
			}

			// log out
			suite.jsonRpcClient.Logout()
		}
	}
}

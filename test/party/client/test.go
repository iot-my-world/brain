package client

import (
	"github.com/stretchr/testify/suite"
	jsonRpcClient "gitlab.com/iotTracker/brain/communication/jsonRpc/client"
	basicJsonRpcClient "gitlab.com/iotTracker/brain/communication/jsonRpc/client/basic"
	authJsonRpcAdaptor "gitlab.com/iotTracker/brain/security/auth/service/adaptor/jsonRpc"
	partyRegistrarJsonRpcAdaptor "gitlab.com/iotTracker/brain/party/registrar/adaptor/jsonRpc"
	"fmt"
	"gitlab.com/iotTracker/brain/party/user"
	"strings"
	"gopkg.in/square/go-jose.v2"
	"reflect"
	"gitlab.com/iotTracker/brain/security/wrappedClaims"
	"gitlab.com/iotTracker/brain/security/claims"
	"encoding/json"
	"gitlab.com/iotTracker/brain/security/claims/registerClientUser"
)

type Client struct {
	suite.Suite
	jsonRpcClient jsonRpcClient.Client
}

func (suite *Client) SetupTest() {
	// create the client
	suite.jsonRpcClient = basicJsonRpcClient.New("http://localhost:9010/api")
}

func (suite *Client) TestClientInviteAndRegisterUsers() {
	for companyOwner := range EntitiesAndAdminUsersToCreate {
		for clientDataEntityIdx := range EntitiesAndAdminUsersToCreate[companyOwner] {
			clientDataEntity := &EntitiesAndAdminUsersToCreate[companyOwner][clientDataEntityIdx]
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
				clientUser := user.User{
					EmailAddress:    (*userEntity).EmailAddress,
					ParentPartyType: suite.jsonRpcClient.Claims().PartyDetails().ParentPartyType,
					ParentId:        suite.jsonRpcClient.Claims().PartyDetails().ParentId,
					PartyType:       suite.jsonRpcClient.Claims().PartyDetails().PartyType,
					PartyId:         suite.jsonRpcClient.Claims().PartyDetails().PartyId,
				}

				// invite the user
				inviteClientUserResponse := partyRegistrarJsonRpcAdaptor.InviteClientUserResponse{}
				if err := suite.jsonRpcClient.JsonRpcRequest(
					"PartyRegistrar.InviteClientUser",
					partyRegistrarJsonRpcAdaptor.InviteClientUserRequest{
						User: clientUser,
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
				registerJsonRpcClient := basicJsonRpcClient.New("http://localhost:9010/api")
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
package client

import (
	"github.com/stretchr/testify/suite"
	jsonRpcClient "gitlab.com/iotTracker/brain/communication/jsonRpc/client"
	basicJsonRpcClient "gitlab.com/iotTracker/brain/communication/jsonRpc/client/basic"
	authJsonRpcAdaptor "gitlab.com/iotTracker/brain/security/auth/service/adaptor/jsonRpc"
	"fmt"
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

			//// invite and register all of the users
			//for userIdx := range (*clientDataEntity).Users {
			//
			//}

			// log out
			suite.jsonRpcClient.Logout()
		}
	}
}

package client

import (
	jsonRpcClient "github.com/iot-my-world/brain/communication/jsonRpc/client"
	basicJsonRpcClient "github.com/iot-my-world/brain/communication/jsonRpc/client/basic"
	"github.com/iot-my-world/brain/party/client"
	clientAdministrator "github.com/iot-my-world/brain/party/client/administrator"
	clientRecordHandler "github.com/iot-my-world/brain/party/client/recordHandler"
	clientJsonRpcAdministrator "github.com/iot-my-world/brain/party/company/administrator/jsonRpc"
	clientJsonRpcRecordHandler "github.com/iot-my-world/brain/party/company/recordHandler/jsonRpc"
	partyRegistrar "github.com/iot-my-world/brain/party/registrar"
	partyJsonRpcRegistrar "github.com/iot-my-world/brain/party/registrar/jsonRpc"
	authJsonRpcAdaptor "github.com/iot-my-world/brain/security/authorization/service/adaptor/jsonRpc"
	humanUser "github.com/iot-my-world/brain/user/human"
	"github.com/stretchr/testify/suite"
)

func New(
	url string,
	user humanUser.User,
	testData *[]Data,
) *test {
	return &test{
		testData:      testData,
		user:          user,
		jsonRpcClient: basicJsonRpcClient.New(url),
	}
}

type test struct {
	suite.Suite
	jsonRpcClient       jsonRpcClient.Client
	clientRecordHandler clientRecordHandler.RecordHandler
	clientAdministrator clientAdministrator.Administrator
	partyRegistrar      partyRegistrar.Registrar
	user                humanUser.User
	testData            *[]Data
}

type Data struct {
	Company   client.Client
	AdminUser humanUser.User
	Users     []humanUser.User
}

func (suite *test) SetupTest() {

	// log in the client
	if err := suite.jsonRpcClient.Login(authJsonRpcAdaptor.LoginRequest{
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
}

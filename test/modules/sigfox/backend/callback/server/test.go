package server

import (
	"fmt"
	jsonRpcClient "github.com/iot-my-world/brain/pkg/api/jsonRpc/client"
	basicJsonRpcClient "github.com/iot-my-world/brain/pkg/api/jsonRpc/client/basic"
	jsonRpcServerAuthenticator "github.com/iot-my-world/brain/pkg/api/jsonRpc/server/authenticator"
	"github.com/iot-my-world/brain/pkg/device/sigbug"
	"github.com/iot-my-world/brain/pkg/search/identifier/name"
	sigfoxBackend "github.com/iot-my-world/brain/pkg/sigfox/backend"
	sigfoxBackendCallbackServer "github.com/iot-my-world/brain/pkg/sigfox/backend/callback/server"
	sigfoxBackendJsonRpcCallbackServer "github.com/iot-my-world/brain/pkg/sigfox/backend/callback/server/jsonRpc"
	sigfoxBackendRecordHandler "github.com/iot-my-world/brain/pkg/sigfox/backend/recordHandler"
	sigfoxBackendJsonRpcRecordHandler "github.com/iot-my-world/brain/pkg/sigfox/backend/recordHandler/jsonRpc"
	humanUser "github.com/iot-my-world/brain/pkg/user/human"
	sigbugGPSTestData "github.com/iot-my-world/brain/test/data/sigbug/gps"
	"github.com/stretchr/testify/suite"
)

func New(
	user humanUser.User,
	humanUserUrl string,
	apiUserUrl string,
	backend sigfoxBackend.Backend,
	testData []Data,
) *test {
	return &test{
		testData:               testData,
		humanUserJsonRpcClient: basicJsonRpcClient.New(humanUserUrl),
		apiUserJsonRpcClient:   basicJsonRpcClient.New(apiUserUrl),
		backend:                backend,
		user:                   user,
	}
}

type test struct {
	suite.Suite
	testData                    []Data
	humanUserJsonRpcClient      jsonRpcClient.Client
	apiUserJsonRpcClient        jsonRpcClient.Client
	backend                     sigfoxBackend.Backend
	user                        humanUser.User
	sigfoxBackendRecordHandler  sigfoxBackendRecordHandler.RecordHandler
	sigfoxBackendCallbackServer sigfoxBackendCallbackServer.Server
}

type Data struct {
	Sigbug  sigbug.Sigbug
	GPSData []sigbugGPSTestData.Data
}

func (suite *test) SetupTest() {
	// log in the human user client
	if err := suite.humanUserJsonRpcClient.Login(jsonRpcServerAuthenticator.LoginRequest{
		UsernameOrEmailAddress: suite.user.Username,
		Password:               string(suite.user.Password),
	}); err != nil {
		suite.FailNow("human user log in error", err.Error())
		return
	}

	// create json rpc record handler
	suite.sigfoxBackendRecordHandler = sigfoxBackendJsonRpcRecordHandler.New(suite.humanUserJsonRpcClient)

	// retrieve the sigfox backend
	retrieveResponse, err := suite.sigfoxBackendRecordHandler.Retrieve(&sigfoxBackendRecordHandler.RetrieveRequest{
		Identifier: name.Identifier{
			Name: suite.backend.Name,
		},
	})
	if err != nil {
		suite.FailNow("error retrieving sigfox backend", err.Error())
		return
	}

	// populate the backend
	suite.backend = retrieveResponse.Backend

	// set token in backend json rpc client
	if err := suite.apiUserJsonRpcClient.SetJWT(suite.backend.Token); err != nil {
		suite.FailNow("error setting token in api user json rpc client", err.Error())
		return
	}

	// create json rpc sigfox backend callback server
	suite.sigfoxBackendCallbackServer = sigfoxBackendJsonRpcCallbackServer.New(suite.apiUserJsonRpcClient)
}

func (suite *test) TestSigfoxBackendCallbackServer1() {
	for testDataIdx := range suite.testData {
		fmt.Println("create!", suite.testData[testDataIdx].GPSData)
	}
}

package server

import (
	jsonRpcClient "github.com/iot-my-world/brain/pkg/api/jsonRpc/client"
	basicJsonRpcClient "github.com/iot-my-world/brain/pkg/api/jsonRpc/client/basic"
	jsonRpcServerAuthenticator "github.com/iot-my-world/brain/pkg/api/jsonRpc/server/authenticator"
	"github.com/iot-my-world/brain/pkg/search/identifier/name"
	sigfoxBackend "github.com/iot-my-world/brain/pkg/sigfox/backend"
	sigfoxBackendRecordHandler "github.com/iot-my-world/brain/pkg/sigfox/backend/recordHandler"
	sigfoxBackendJsonRpcRecordHandler "github.com/iot-my-world/brain/pkg/sigfox/backend/recordHandler/jsonRpc"
	humanUser "github.com/iot-my-world/brain/pkg/user/human"
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
	testData                   []Data
	humanUserJsonRpcClient     jsonRpcClient.Client
	apiUserJsonRpcClient       jsonRpcClient.Client
	backend                    sigfoxBackend.Backend
	user                       humanUser.User
	sigfoxBackendRecordHandler sigfoxBackendRecordHandler.RecordHandler
}

type Data struct {
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
}

func (suite *test) TestSigfoxBackendCallbackServer1() {
	suite.T().Log("awer")
}

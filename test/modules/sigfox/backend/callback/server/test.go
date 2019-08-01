package server

import (
	jsonRpcClient "github.com/iot-my-world/brain/pkg/api/jsonRpc/client"
	basicJsonRpcClient "github.com/iot-my-world/brain/pkg/api/jsonRpc/client/basic"
	"github.com/iot-my-world/brain/pkg/security/claims/sigfoxBackend"
	"github.com/stretchr/testify/suite"
)

func New(
	humanUserUrl string,
	apiUserUrl string,
	backend sigfoxBackend.SigfoxBackend,
	testData []Data,
) *test {
	return &test{
		testData:               testData,
		humanUserJsonRpcClient: basicJsonRpcClient.New(humanUserUrl),
		apiUserJsonRpcClient:   basicJsonRpcClient.New(apiUserUrl),
		backend:                backend,
	}
}

type test struct {
	suite.Suite
	testData               []Data
	humanUserJsonRpcClient jsonRpcClient.Client
	apiUserJsonRpcClient   jsonRpcClient.Client
	backend                sigfoxBackend.SigfoxBackend
}

type Data struct {
}

func (suite *test) SetupTest() {

}

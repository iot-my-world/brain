package server

import (
	jsonRpcClient "github.com/iot-my-world/brain/pkg/api/jsonRpc/client"
	basicJsonRpcClient "github.com/iot-my-world/brain/pkg/api/jsonRpc/client/basic"
	"github.com/stretchr/testify/suite"
)

func New(
	url string,
	testData []Data,
) *test {
	return &test{
		testData:      testData,
		jsonRpcClient: basicJsonRpcClient.New(url),
	}
}

type test struct {
	suite.Suite
	jsonRpcClient jsonRpcClient.Client
	testData      []Data
}

type Data struct {
}

func (suite *test) SetupTest() {

}

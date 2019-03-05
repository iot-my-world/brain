package main

import (
	"testing"
	"github.com/stretchr/testify/suite"
	jsonRpcClient "gitlab.com/iotTracker/brain/communication/jsonRpc/client"
	basicJsonRpcClient "gitlab.com/iotTracker/brain/communication/jsonRpc/client/basic"
	authJsonRpcAdaptor "gitlab.com/iotTracker/brain/security/auth/service/adaptor/jsonRpc"
	"fmt"
)

// Define the suite, and absorb the built-in basic suite
// functionality from testify - including a T() method which
// returns the current testing context
type ExampleTestSuite struct {
	suite.Suite
	VariableThatShouldStartAtFive int
	jsonRpcClient                 jsonRpcClient.Client
}

// Make sure that VariableThatShouldStartAtFive is set to five
// before each test
func (suite *ExampleTestSuite) SetupTest() {
	suite.jsonRpcClient = basicJsonRpcClient.New("http://localhost:9010/api")
}

// All methods that begin with "Test" are run as tests within a
// suite.
func (suite *ExampleTestSuite) TestExample() {
	loginRequest := authJsonRpcAdaptor.LoginRequest{
		UsernameOrEmailAddress: "root",
		Password:               "12345",
	}

	jsonRpcRequest := jsonRpcClient.NewRequest(
		"1234",
			"Auth.Login",
		)
	jsonRpcRequest.Params = [1]interface{}{loginRequest}

	response, err := suite.jsonRpcClient.Post(&jsonRpcRequest)
	if err != nil {
		suite.T().Errorf(err.Error())
	}

	fmt.Println("success!", response)
}

// In order for 'go test' to run this suite, we need to create
// a normal test function and pass our suite to suite.Run
func TestBrain(t *testing.T) {
	suite.Run(t, new(ExampleTestSuite))
}

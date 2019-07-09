package sigbug

import (
	jsonRpcClient "github.com/iot-my-world/brain/pkg/communication/jsonRpc/client"
	basicJsonRpcClient "github.com/iot-my-world/brain/pkg/communication/jsonRpc/client/basic"
	"github.com/iot-my-world/brain/pkg/device/sigbug"
	sigugGPSReadings "github.com/iot-my-world/brain/pkg/device/sigbug/reading/gps"
	authorizationAdministrator "github.com/iot-my-world/brain/pkg/security/authorization/administrator"
	humanUser "github.com/iot-my-world/brain/pkg/user/human"
	"github.com/stretchr/testify/suite"
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
	jsonRpcClient jsonRpcClient.Client
	user          humanUser.User
	testData      []Data
}

type Data struct {
	Device      sigbug.Sigbug
	GPSReadings []sigugGPSReadings.Reading
}

func (suite *test) SetupTest() {

	// log in the client
	if err := suite.jsonRpcClient.Login(authorizationAdministrator.LoginRequest{
		UsernameOrEmailAddress: suite.user.Username,
		Password:               string(suite.user.Password),
	}); err != nil {
		suite.Fail("log in error", err.Error())
		return
	}

	// set up service provider clients that use jsonRpcClient
}

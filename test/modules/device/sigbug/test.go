package sigbug

import (
	jsonRpcClient "github.com/iot-my-world/brain/pkg/communication/jsonRpc/client"
	basicJsonRpcClient "github.com/iot-my-world/brain/pkg/communication/jsonRpc/client/basic"
	"github.com/iot-my-world/brain/pkg/device/sigbug"
	sigbugAdministrator "github.com/iot-my-world/brain/pkg/device/sigbug/administrator"
	sigbugJsonRpcAdministrator "github.com/iot-my-world/brain/pkg/device/sigbug/administrator/jsonRpc"
	sigbugGPSReadings "github.com/iot-my-world/brain/pkg/device/sigbug/reading/gps"
	sigbugRecordHandler "github.com/iot-my-world/brain/pkg/device/sigbug/recordHandler"
	sigbugJsonRpcRecordHandler "github.com/iot-my-world/brain/pkg/device/sigbug/recordHandler/jsonRpc"
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
	jsonRpcClient       jsonRpcClient.Client
	user                humanUser.User
	testData            []Data
	sigbugAdministrator sigbugAdministrator.Administrator
	sigbugRecordHandler sigbugRecordHandler.RecordHandler
}

type Data struct {
	Device      sigbug.Sigbug
	GPSReadings []sigbugGPSReadings.Reading
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
	suite.sigbugAdministrator = sigbugJsonRpcAdministrator.New(suite.jsonRpcClient)
	suite.sigbugRecordHandler = sigbugJsonRpcRecordHandler.New(suite.jsonRpcClient)
}

func (suite *test) TestSigbug1Create() {
	// create all sigbugs in test data
	for _, data := range suite.testData {
		// if owner party name set, retiev
	}
}

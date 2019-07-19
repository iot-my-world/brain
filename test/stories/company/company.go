package company

import (
	basicJsonRpcClient "github.com/iot-my-world/brain/pkg/api/jsonRpc/client/basic"
	jsonRpcServerAuthenticator "github.com/iot-my-world/brain/pkg/api/jsonRpc/server/authenticator"
	"github.com/iot-my-world/brain/test/data/environment"
	sigbugDeviceTestModule "github.com/iot-my-world/brain/test/modules/device/sigbug"
	clientTestModule "github.com/iot-my-world/brain/test/modules/party/client"
	clientStoryTestData "github.com/iot-my-world/brain/test/stories/client/data"
	companyStoryTestData "github.com/iot-my-world/brain/test/stories/company/data"
	systemStoryTestData "github.com/iot-my-world/brain/test/stories/system/data"
	"github.com/stretchr/testify/suite"
)

func New() *test {
	return &test{}
}

type test struct {
	suite.Suite
}

func (t *test) SetupTest() {

}

func (t *test) TestCompany() {
	// log in a json rpc client
	jsonRpcClient := basicJsonRpcClient.New(environment.APIUserURL)
	if err := jsonRpcClient.Login(jsonRpcServerAuthenticator.LoginRequest{
		UsernameOrEmailAddress: systemStoryTestData.User.Username,
		Password:               string(systemStoryTestData.User.Password),
	}); err != nil {
		t.Fail("log in error", err.Error())
		return
	}

	for _, companyData := range companyStoryTestData.TestData {
		// get client data for client owned by this company
		clientData, found := clientStoryTestData.TestData[companyData.CompanyTestData.Company.Name]
		if !found {
			t.FailNow("no client data for company")
			return
		}

		// build client test data and run client tests
		clientTestData := make([]clientTestModule.Data, 0)
		for _, cData := range clientData {
			clientTestData = append(clientTestData, cData.ClientTestData)
		}
		suite.Run(t.T(), clientTestModule.New(
			environment.BrainHumanUserURL,
			companyData.CompanyTestData.AdminUser,
			clientTestData,
		))

		// build sigbug test data and run sigbug tests
		sigbugTestData := make([]sigbugDeviceTestModule.Data, 0)
		for _, sigbugDevice := range companyData.SigbugDevices {
			sigbugTestData = append(sigbugTestData, sigbugDeviceTestModule.Data{
				Device:      sigbugDevice,
				GPSReadings: nil,
			})
		}
		suite.Run(t.T(), sigbugDeviceTestModule.New(
			environment.BrainHumanUserURL,
			environment.APIUserURL,
			companyData.CompanyTestData.AdminUser,
			sigbugTestData,
			"123",
		))
	}
}

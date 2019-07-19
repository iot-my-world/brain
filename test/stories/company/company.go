package company

import (
	basicJsonRpcClient "github.com/iot-my-world/brain/pkg/api/jsonRpc/client/basic"
	jsonRpcServerAuthenticator "github.com/iot-my-world/brain/pkg/api/jsonRpc/server/authenticator"
	"github.com/iot-my-world/brain/pkg/search/identifier/name"
	sigfoxBackendRecordHandler "github.com/iot-my-world/brain/pkg/sigfox/backend/recordHandler"
	sigfoxBackendJsonRpcRecordHandler "github.com/iot-my-world/brain/pkg/sigfox/backend/recordHandler/jsonRpc"
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
	jsonRpcClient := basicJsonRpcClient.New(environment.BrainHumanUserURL)
	if err := jsonRpcClient.Login(jsonRpcServerAuthenticator.LoginRequest{
		UsernameOrEmailAddress: systemStoryTestData.User.Username,
		Password:               string(systemStoryTestData.User.Password),
	}); err != nil {
		t.FailNow("log in error", err.Error())
		return
	}
	backendRecordHandler := sigfoxBackendJsonRpcRecordHandler.New(jsonRpcClient)

	// retrieve sigfox backend
	backendRetrieveResponse, err := backendRecordHandler.Retrieve(&sigfoxBackendRecordHandler.RetrieveRequest{
		Identifier: name.Identifier{
			Name: systemStoryTestData.SigfoxBackendTestData[0].Backend.Name,
		},
	})
	if err != nil {
		t.FailNow("error retrieving sigfox backend", err.Error())
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
			backendRetrieveResponse.Backend.Token,
		))
	}
}

package system

import (
	"github.com/iot-my-world/brain/pkg/workbook"
	"github.com/iot-my-world/brain/test/data/environment"
	clientTestModule "github.com/iot-my-world/brain/test/modules/party/client"
	companyTestModule "github.com/iot-my-world/brain/test/modules/party/company"
	sigfoxBackendTestModule "github.com/iot-my-world/brain/test/modules/sigfox/backend"
	sigfoxBackendCallbackServerTestModule "github.com/iot-my-world/brain/test/modules/sigfox/backend/callback/server"
	clientStoryTestData "github.com/iot-my-world/brain/test/stories/client/data"
	companyTestStoryData "github.com/iot-my-world/brain/test/stories/company/data"
	systemTestStoryData "github.com/iot-my-world/brain/test/stories/system/data"
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

func (t *test) TestSystem() {
	// perform system company tests
	companyTestData := make([]companyTestModule.Data, 0)
	for _, companyData := range companyTestStoryData.TestData {
		companyTestData = append(companyTestData, companyData.CompanyTestData)
	}
	suite.Run(t.T(), companyTestModule.New(
		environment.BrainHumanUserURL,
		systemTestStoryData.User,
		companyTestData,
	))

	// perform system client tests
	clientData, found := clientStoryTestData.TestData["root"]
	if !found {
		t.FailNow("root client data not found")
		return
	}

	clientTestData := make([]clientTestModule.Data, 0)
	for _, clientData := range clientData {
		clientTestData = append(clientTestData, clientData.ClientTestData)
	}
	suite.Run(t.T(), clientTestModule.New(
		environment.BrainHumanUserURL,
		systemTestStoryData.User,
		clientTestData,
	))

	for _, sigfoxBackendData := range systemTestStoryData.SigfoxBackendTestData {
		// create, update, retrieve etc.
		suite.Run(t.T(), sigfoxBackendTestModule.New(
			environment.BrainHumanUserURL,
			systemTestStoryData.User,
			[]sigfoxBackendTestModule.Data{
				sigfoxBackendData,
			},
		))

		// parse test data
		gpsDataWorkbook, err := workbook.New()

		// tests logged in as backend
		suite.Run(t.T(), sigfoxBackendCallbackServerTestModule.New(
			systemTestStoryData.User,
			environment.BrainHumanUserURL,
			environment.APIUserURL,
			sigfoxBackendData.Backend,
			[]sigfoxBackendCallbackServerTestModule.Data{},
		))
	}

}

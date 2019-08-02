package system

import (
	"github.com/iot-my-world/brain/test/data/environment"
	sigbugGPSTestData "github.com/iot-my-world/brain/test/data/sigbug/gps"
	sigbugGPSTestDataGenerator "github.com/iot-my-world/brain/test/data/sigbug/gps/generator"
	clientTestModule "github.com/iot-my-world/brain/test/modules/party/client"
	companyTestModule "github.com/iot-my-world/brain/test/modules/party/company"
	sigfoxBackendTestModule "github.com/iot-my-world/brain/test/modules/sigfox/backend"
	sigfoxBackendCallbackServerTestModule "github.com/iot-my-world/brain/test/modules/sigfox/backend/callback/server"
	clientStoryTestData "github.com/iot-my-world/brain/test/stories/client/data"
	companyTestStoryData "github.com/iot-my-world/brain/test/stories/company/data"
	systemTestStoryData "github.com/iot-my-world/brain/test/stories/system/data"
	"github.com/stretchr/testify/suite"
	"math"
)

func New() *test {
	return &test{}
}

type test struct {
	suite.Suite
}

func (t *test) SetupTest() {

}

const noGPSReadingsToTake = 10

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
		gpsDataMap, err := sigbugGPSTestDataGenerator.Generate()
		if err != nil {
			t.FailNow("error getting sigbug gps test data", err)
			return
		}

		// get 10 readings from each test journey data set
		testGPSData := make([]sigbugGPSTestData.Data, 0)
		for journeyName := range gpsDataMap {
			if noGPSReadingsToTake > len(gpsDataMap[journeyName]) {
				// if the number to be taken is greater than the size of the set
				// then take all
				testGPSData = append(testGPSData, gpsDataMap[journeyName]...)
				continue
			}
			// otherwise sample the set
			for i := 0; i < noGPSReadingsToTake; i++ {
				sampleIdx := int(math.Ceil(float64(i*len(gpsDataMap[journeyName])) / float64(noGPSReadingsToTake)))
				if sampleIdx < 0 || sampleIdx == len(gpsDataMap[journeyName]) {
					t.FailNow("sample index invalid", sampleIdx)
					return
				}
				testGPSData = append(
					testGPSData,
					gpsDataMap[journeyName][sampleIdx],
				)
			}
		}

		// tests logged in as backend
		suite.Run(t.T(), sigfoxBackendCallbackServerTestModule.New(
			systemTestStoryData.User,
			environment.BrainHumanUserURL,
			environment.APIUserURL,
			sigfoxBackendData.Backend,
			[]sigfoxBackendCallbackServerTestModule.Data{
				{
					Sigbug:  clientData[0].SigbugDevices[0],
					GPSData: testGPSData,
				},
			},
		))
	}

}

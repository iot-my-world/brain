package client

import (
	"github.com/iot-my-world/brain/test/data/environment"
	sigbugDeviceTestModule "github.com/iot-my-world/brain/test/modules/device/sigbug"
	clientStoryTestData "github.com/iot-my-world/brain/test/stories/client/data"
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

func (t *test) TestClient() {
	for _, clientsTestData := range clientStoryTestData.TestData {
		for _, data := range clientsTestData {
			// build sigbug test data and run sigbug tests
			sigbugTestData := make([]sigbugDeviceTestModule.Data, 0)
			for _, sigbugDevice := range data.SigbugDevices {
				sigbugTestData = append(sigbugTestData, sigbugDeviceTestModule.Data{
					Device:      sigbugDevice,
					GPSReadings: nil,
				})
			}
			suite.Run(t.T(), sigbugDeviceTestModule.New(
				environment.BrainHumanUserURL,
				environment.APIUserURL,
				data.ClientTestData.AdminUser,
				sigbugTestData,
				"",
			))
		}
	}
}

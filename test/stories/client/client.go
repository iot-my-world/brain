package client

import (
	"github.com/iot-my-world/brain/test/data"
	sigbugDeviceTestModule "github.com/iot-my-world/brain/test/modules/device/sigbug"
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
	for _, clientsTestData := range TestData {
		for _, clientTestData := range clientsTestData {
			// build sigbug test data and run sigbug tests
			sigbugTestData := make([]sigbugDeviceTestModule.Data, 0)
			for _, sigbugDevice := range clientTestData.SigbugDevices {
				sigbugTestData = append(sigbugTestData, sigbugDeviceTestModule.Data{
					Device:      sigbugDevice,
					GPSReadings: nil,
				})
			}
			suite.Run(t.T(), sigbugDeviceTestModule.New(
				data.BrainURL,
				clientTestData.ClientTestData.AdminUser,
				sigbugTestData,
			))
		}
	}
}

package company

import (
	"github.com/iot-my-world/brain/test/data"
	sigbugDeviceTestModule "github.com/iot-my-world/brain/test/modules/device/sigbug"
	clientTestModule "github.com/iot-my-world/brain/test/modules/party/client"
	"github.com/iot-my-world/brain/test/stories/client"
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
	for _, companyData := range TestData {
		// get client data for client owned by this company
		clientData, found := client.TestData[companyData.CompanyTestData.Company.Name]
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
			data.BrainURL,
			companyData.CompanyTestData.AdminUser,
			clientTestData,
		))

		// build sigbug test data and run sigbug tests
		sigbugTestData := make([]sigbugDeviceTestModule.Data, 0)
		for _, cData := range clientData {
			for _, sigbugDevice := range cData.SigbugDevices {
				sigbugTestData = append(sigbugTestData, sigbugDeviceTestModule.Data{
					Device:      sigbugDevice,
					GPSReadings: nil,
				})
			}
		}
		suite.Run(t.T(), sigbugDeviceTestModule.New(
			data.BrainURL,
			companyData.CompanyTestData.AdminUser,
			sigbugTestData,
		))
	}
}

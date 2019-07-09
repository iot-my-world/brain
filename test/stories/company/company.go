package company

import (
	"github.com/iot-my-world/brain/test/data"
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
		// build client data
		clientTestData := make([]clientTestModule.Data, 0)
		for _, data := range clientData {
			clientTestData = append(clientTestData, data.ClientTestData)
		}
		suite.Run(t.T(), clientTestModule.New(
			data.BrainURL,
			companyData.CompanyTestData.AdminUser,
			clientTestData,
		))
	}
}

package system

import (
	"github.com/iot-my-world/brain/test/data"
	clientTestModule "github.com/iot-my-world/brain/test/modules/party/client"
	companyTestModule "github.com/iot-my-world/brain/test/modules/party/company"
	"github.com/iot-my-world/brain/test/stories/client"
	"github.com/iot-my-world/brain/test/stories/company"
	"github.com/stretchr/testify/suite"
	"testing"
)

func Test(t *testing.T) {
	// perform system company tests
	companyTestData := make([]companyTestModule.Data, 0)
	for _, companyData := range company.TestData {
		companyTestData = append(companyTestData, companyData.CompanyTestData)
	}
	suite.Run(t, companyTestModule.New(
		data.BrainURL,
		User,
		companyTestData,
	))

	// perform system client tests
	clientData, found := client.TestData["root"]
	if !found {
		t.Logf("root client data not found")
		t.FailNow()
		return
	}

	clientTestData := make([]clientTestModule.Data, 0)
	for _, clientData := range clientData {
		clientTestData = append(clientTestData, clientData.ClientTestData)
	}
	suite.Run(t, clientTestModule.New(
		data.BrainURL,
		User,
		clientTestData,
	))
}

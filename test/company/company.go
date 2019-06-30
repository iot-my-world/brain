package company

import (
	"github.com/iot-my-world/brain/test/data/environment"
	clientTestData "github.com/iot-my-world/brain/test/data/party/client"
	companyTestData "github.com/iot-my-world/brain/test/data/party/company"
	clientTest "github.com/iot-my-world/brain/test/modules/party/client"
	"github.com/stretchr/testify/suite"
	"testing"
)

func Test(t *testing.T) {
	for _, companyData := range companyTestData.TestData {
		clientData, found := clientTestData.TestData[companyData.Company.Name]
		if !found {
			t.Fatalf("no client data for company")
			return
		}
		suite.Run(t, clientTest.New(
			environment.BrainURL,
			companyData.AdminUser,
			clientData,
		))
	}
}

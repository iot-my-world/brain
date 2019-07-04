package company

import (
	"github.com/iot-my-world/brain/test/data"
	clientTestModule "github.com/iot-my-world/brain/test/modules/party/client"
	"github.com/iot-my-world/brain/test/stories/client"
	"github.com/stretchr/testify/suite"
	"testing"
)

func Test(t *testing.T) {
	for _, companyData := range TestData {
		clientData, found := client.TestData[companyData.Company.Name]
		if !found {
			t.Fatalf("no client data for company")
			return
		}
		suite.Run(t, clientTestModule.New(
			data.BrainURL,
			companyData.AdminUser,
			clientData,
		))
	}
}

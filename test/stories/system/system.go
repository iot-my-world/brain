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
	suite.Run(t, companyTestModule.New(
		data.BrainURL,
		User,
		company.TestData,
	))

	// perform system client tests
	clientData, found := client.TestData["root"]
	if !found {
		t.Logf("root client data not found")
		t.FailNow()
		return
	}
	suite.Run(t, clientTestModule.New(
		data.BrainURL,
		User,
		clientData,
	))
}

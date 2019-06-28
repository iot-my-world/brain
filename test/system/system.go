package system

import (
	"github.com/iot-my-world/brain/test/data/environment"
	companyTestData "github.com/iot-my-world/brain/test/data/party/company"
	systemTestData "github.com/iot-my-world/brain/test/data/party/system"
	companyTest "github.com/iot-my-world/brain/test/modules/party/company"
	"github.com/stretchr/testify/suite"
	"testing"
)

func Test(t *testing.T) {
	// perform system company tests
	suite.Run(t, companyTest.New(
		environment.BrainURL,
		systemTestData.User,
		companyTestData.TestData,
	))
}

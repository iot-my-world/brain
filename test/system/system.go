package system

import (
	testData "github.com/iot-my-world/brain/test/data"
	companyTest "github.com/iot-my-world/brain/test/modules/party/company"
	"github.com/stretchr/testify/suite"
	"testing"
)

func Test(t *testing.T) {
	// perform system company tests
	suite.Run(t, companyTest.New(
		testData.BrainURL,
		testData.SystemUser,
		&testData.CompanyTestData,
	))
}

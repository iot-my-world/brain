package system

import (
	"github.com/iot-my-world/brain/test"
	companyTest "github.com/iot-my-world/brain/test/modules/party/company"
	"github.com/iot-my-world/brain/test/stories/company"
	"github.com/stretchr/testify/suite"
	"testing"
)

func Test(t *testing.T) {
	// perform system company tests
	suite.Run(t, companyTest.New(
		test.BrainURL,
		User,
		company.TestData,
	))
}

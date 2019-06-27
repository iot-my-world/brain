package company

import (
	"github.com/iot-my-world/brain/test/data/environment"
	clientTestData "github.com/iot-my-world/brain/test/data/party/client"
	systemTestData "github.com/iot-my-world/brain/test/data/party/system"
	clientTest "github.com/iot-my-world/brain/test/modules/party/client"
	"github.com/stretchr/testify/suite"
	"testing"
)

func Test(t *testing.T) {
	// perform company client tests
	suite.Run(t, clientTest.New(
		environment.BrainURL,
		systemTestData.User,
		clientTestData.TestData,
	))
}

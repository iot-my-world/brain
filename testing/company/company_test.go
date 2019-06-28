package company

import (
	clientTest "github.com/iot-my-world/brain/testing/company/client"
	userTest "github.com/iot-my-world/brain/testing/company/user"
	"github.com/stretchr/testify/suite"
	"testing"
)

func TestCompany(t *testing.T) {
	suite.Run(t, new(userTest.User))
	suite.Run(t, new(clientTest.Client))
}

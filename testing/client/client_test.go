package client

import (
	userTest "github.com/iot-my-world/brain/testing/client/user"
	"github.com/stretchr/testify/suite"
	"testing"
)

func TestCompany(t *testing.T) {
	suite.Run(t, new(userTest.User))
}

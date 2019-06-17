package client

import (
	"github.com/stretchr/testify/suite"
	userTest "gitlab.com/iotTracker/brain/test/client/user"
	"testing"
)

func TestCompany(t *testing.T) {
	suite.Run(t, new(userTest.User))
}

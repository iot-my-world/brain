package company

import (
	"github.com/stretchr/testify/suite"
	clientTest "gitlab.com/iotTracker/brain/test/company/client"
	userTest "gitlab.com/iotTracker/brain/test/company/user"
	"testing"
)

func TestCompany(t *testing.T) {
	suite.Run(t, new(userTest.User))
	suite.Run(t, new(clientTest.Client))
}

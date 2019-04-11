package system

import (
	"github.com/stretchr/testify/suite"
	companyTest "gitlab.com/iotTracker/brain/test/system/company"
	"testing"
)

func TestSystem(t *testing.T) {
	suite.Run(t, new(companyTest.Company))
}
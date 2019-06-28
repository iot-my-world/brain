package system

import (
	companyTest "github.com/iot-my-world/brain/testing/system/company"
	"github.com/stretchr/testify/suite"
	"testing"
)

func TestSystem(t *testing.T) {
	suite.Run(t, new(companyTest.Company))
}

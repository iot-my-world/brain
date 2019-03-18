package device

import (
	"github.com/stretchr/testify/suite"
	systemTest "gitlab.com/iotTracker/brain/test/device/system"
	"testing"
)

func TestDevice(t *testing.T) {
	suite.Run(t, new(systemTest.System))
}

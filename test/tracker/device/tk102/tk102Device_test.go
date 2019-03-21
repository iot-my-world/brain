package tk102

import (
	"github.com/stretchr/testify/suite"
	systemTest "gitlab.com/iotTracker/brain/test/tracker/device/tk102/system"
	"testing"
)

func TestDevice(t *testing.T) {
	suite.Run(t, new(systemTest.System))
}

package reading

import (
	"github.com/stretchr/testify/suite"
	tk102DeviceSytemTest "gitlab.com/iotTracker/brain/test/tracker/device/tk102/system"
	readingSystemTest "gitlab.com/iotTracker/brain/test/tracker/reading/system"
	"testing"
)

func TestTracker(t *testing.T) {
	suite.Run(t, new(tk102DeviceSytemTest.System))
	suite.Run(t, new(readingSystemTest.System))
}

package sf001

import (
	systemTest "github.com/iot-my-world/brain/test/tracker/sf001/system"
	"github.com/stretchr/testify/suite"
	"testing"
)

func TestDevice(t *testing.T) {
	suite.Run(t, new(systemTest.System))
}

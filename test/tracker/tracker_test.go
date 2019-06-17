package reading

import (
	sf001TrackerSystemTest "github.com/iot-my-world/brain/test/tracker/sf001/system"
	"github.com/stretchr/testify/suite"
	"testing"
)

func TestTracker(t *testing.T) {
	suite.Run(t, new(sf001TrackerSystemTest.System))
}

package reading

import (
	"github.com/stretchr/testify/suite"
	sf001TrackerSystemTest "gitlab.com/iotTracker/brain/test/tracker/sf001/system"
	"testing"
)

func TestTracker(t *testing.T) {
	suite.Run(t, new(sf001TrackerSystemTest.System))
}

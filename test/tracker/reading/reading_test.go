package reading

import (
	"github.com/stretchr/testify/suite"
	readingSystemTest "gitlab.com/iotTracker/brain/test/tracker/reading/system"
	"testing"
)

func TestReading(t *testing.T) {
	suite.Run(t, new(readingSystemTest.System))
}

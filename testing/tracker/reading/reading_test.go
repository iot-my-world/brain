package reading

import (
	readingSystemTest "github.com/iot-my-world/brain/testing/tracker/reading/system"
	"github.com/stretchr/testify/suite"
	"testing"
)

func TestReading(t *testing.T) {
	suite.Run(t, new(readingSystemTest.System))
}

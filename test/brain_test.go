package main

import (
	"testing"
	"github.com/stretchr/testify/suite"
	systemTest "gitlab.com/iotTracker/brain/test/party/system"
)

// In order for 'go test' to run this suite, we need to create
// a normal test function and pass our suite to suite.Run
func TestBrain(t *testing.T) {
	suite.Run(t, new(systemTest.System))
}

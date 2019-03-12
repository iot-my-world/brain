package main

import (
	"testing"
	"github.com/stretchr/testify/suite"
	systemTest "gitlab.com/iotTracker/brain/test/party/system"
	companyTest "gitlab.com/iotTracker/brain/test/party/company"
	clientTest "gitlab.com/iotTracker/brain/test/party/client"
)

// In order for 'go test' to run this suite, we need to create
// a normal test function and pass our suite to suite.Run
func TestBrain(t *testing.T) {
	suite.Run(t, new(systemTest.System))
	suite.Run(t, new(companyTest.Company))
	suite.Run(t, new(clientTest.Client))
}
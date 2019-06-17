package main

import (
	"github.com/stretchr/testify/suite"
	clientUserTest "gitlab.com/iotTracker/brain/test/client/user"
	companyClientTest "gitlab.com/iotTracker/brain/test/company/client"
	companyUserTest "gitlab.com/iotTracker/brain/test/company/user"
	systemCompanyTest "gitlab.com/iotTracker/brain/test/system/company"
	sf001TrackerSystemTest "gitlab.com/iotTracker/brain/test/tracker/sf001/system"
	"testing"
)

// In order for 'go test' to run this suite, we need to create
// a normal test function and pass our suite to suite.Run
func TestBrain(t *testing.T) {
	// System Tests
	suite.Run(t, new(systemCompanyTest.Company))

	// Company Tests
	suite.Run(t, new(companyUserTest.User))
	suite.Run(t, new(companyClientTest.Client))

	// Client Tests
	suite.Run(t, new(clientUserTest.User))

	// Device Tests
	suite.Run(t, new(sf001TrackerSystemTest.System))
}

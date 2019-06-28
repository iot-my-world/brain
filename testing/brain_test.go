package main

import (
	clientUserTest "github.com/iot-my-world/brain/testing/client/user"
	companyClientTest "github.com/iot-my-world/brain/testing/company/client"
	companyUserTest "github.com/iot-my-world/brain/testing/company/user"
	systemCompanyTest "github.com/iot-my-world/brain/testing/system/company"
	sf001TrackerSystemTest "github.com/iot-my-world/brain/testing/tracker/sf001/system"
	"github.com/stretchr/testify/suite"
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

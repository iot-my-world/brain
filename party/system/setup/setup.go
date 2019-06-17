package setup

import (
	"gitlab.com/iotTracker/brain/party"
	partyRegistrar "gitlab.com/iotTracker/brain/party/registrar"
	partyRegistrarException "gitlab.com/iotTracker/brain/party/registrar/exception"
	"gitlab.com/iotTracker/brain/party/system"
	systemRecordHandler "gitlab.com/iotTracker/brain/party/system/recordHandler"
	systemRecordHandlerException "gitlab.com/iotTracker/brain/party/system/recordHandler/exception"
	systemSetupException "gitlab.com/iotTracker/brain/party/system/setup/exception"
	"gitlab.com/iotTracker/brain/search/identifier/id"
	"gitlab.com/iotTracker/brain/search/identifier/name"
	humanUserLoginClaims "gitlab.com/iotTracker/brain/security/claims/login/user/human"
	humanUser "gitlab.com/iotTracker/brain/user/human"
	"io/ioutil"
	"os"
	"strings"
	"time"
)

var systemEntity = system.System{
	//Id
	Name:              "root",
	AdminEmailAddress: "root@root.com",
}

var systemAdminUser = humanUser.User{
	Name:    "root",
	Surname: "root",

	Username:        "root",
	EmailAddress:    "root@root.com",
	Password:        []byte("12345"),
	Roles:           []string{"root"},
	ParentPartyType: party.System,
	// ParentId: // to be set after creating user
	PartyType: party.System,
	// PartyId:  // to be set after creating user
}

func consumePasswordFile(location string) ([]byte, error) {
	if _, err := os.Stat(location); err != nil {
		return nil, err
	}
	// read the file
	data, err := ioutil.ReadFile(location)
	if err != nil {
		return nil, err
	}
	// remove the file
	if err := os.Remove(location); err != nil {
		return nil, err
	}
	// return the data
	return data, nil
}

func InitialSetup(
	handler systemRecordHandler.RecordHandler,
	registrar partyRegistrar.Registrar,
	rootPasswordLocation string,
	systemClaims *humanUserLoginClaims.Login,
) error {
	// try and retrieve the root system entity
	var systemEntityCreatedOrRetrieved system.System
	systemEntityRetrieveResponse, err := handler.Retrieve(&systemRecordHandler.RetrieveRequest{
		Claims:     systemClaims,
		Identifier: name.Identifier{Name: systemEntity.Name},
	})
	switch err.(type) {
	case nil:
		// this means the system entity already exists
		systemEntityCreatedOrRetrieved = systemEntityRetrieveResponse.System

	case systemRecordHandlerException.NotFound:
		// this means that system must be created now

		// get the password from file if a path is provided
		if rootPasswordLocation != "" {
			pwd, err := consumePasswordFile(rootPasswordLocation)
			if err != nil {
				return systemSetupException.InitialSetup{Reasons: []string{"consume password error", err.Error()}}
			}
			systemAdminUser.Password = []byte(strings.TrimSuffix(string(pwd), "\n"))
		}

		// now try create the system
		systemEntityCreateResponse, err := handler.Create(&systemRecordHandler.CreateRequest{
			System: systemEntity,
		})
		if err != nil {
			return systemSetupException.InitialSetup{Reasons: []string{"create error", err.Error()}}
		}
		systemEntityCreatedOrRetrieved = systemEntityCreateResponse.System

	default:
		// some other error
		return systemSetupException.InitialSetup{Reasons: []string{"retrieve error", err.Error()}}
	}

	// assign the id for the system admin user
	systemAdminUser.PartyId = id.Identifier{Id: systemEntityCreatedOrRetrieved.Id}
	systemAdminUser.ParentId = id.Identifier{Id: systemEntityCreatedOrRetrieved.Id}

	// try and register the system admin user
	registerSystemAdminUserResponse, err := registrar.RegisterSystemAdminUser(&partyRegistrar.RegisterSystemAdminUserRequest{
		Claims: systemClaims,
		User:   systemAdminUser,
	})
	if err != nil {
		switch err.(type) {
		case partyRegistrarException.AlreadyRegistered:
			// this is fine, no issues
		default:
			// something went wrong
			return systemSetupException.InitialSetup{Reasons: []string{"registration error", err.Error()}}
		}
	}

	// set up the system claims
	systemClaims.UserId = id.Identifier{Id: registerSystemAdminUserResponse.User.Id}
	systemClaims.IssueTime = time.Now().Unix()
	// systemClaims.ExpirationTime = ?
	systemClaims.ParentPartyType = party.System
	systemClaims.PartyId = id.Identifier{Id: systemEntityCreatedOrRetrieved.Id}
	systemClaims.PartyType = party.System
	systemClaims.PartyId = id.Identifier{Id: systemEntityCreatedOrRetrieved.Id}

	return nil
}

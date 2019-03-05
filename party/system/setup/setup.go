package setup

import (
	systemRecordHandler "gitlab.com/iotTracker/brain/party/system/recordHandler"
	systemRecordHandlerException "gitlab.com/iotTracker/brain/party/system/recordHandler/exception"
	systemSetupException "gitlab.com/iotTracker/brain/party/system/setup/exception"
	partyRegistrar "gitlab.com/iotTracker/brain/party/registrar"
	partyRegistrarException "gitlab.com/iotTracker/brain/party/registrar/exception"
	"gitlab.com/iotTracker/brain/party/system"
	"gitlab.com/iotTracker/brain/search/identifier/name"
	"os"
	"io/ioutil"
	"gitlab.com/iotTracker/brain/party/user"
	"gitlab.com/iotTracker/brain/party"
	"gitlab.com/iotTracker/brain/search/identifier/id"
	"strings"
)

var systemEntity = system.System{
	//Id
	Name:              "root",
	AdminEmailAddress: "root@root.com",
}

var systemAdminUser = user.User{
	Name:    "root",
	Surname: "root",

	Username:     "root",
	EmailAddress: "root@root.com",
	// Password: // set during system user registration
	Roles:           []string{"root"},
	ParentPartyType: party.System,
	// ParentId: // to be set after creating user
	PartyType: party.System,
	// PartyId:  // to be set after creating user
}

var defaultSystemPassword = "12345"

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

func InitialSetup(handler systemRecordHandler.RecordHandler, registrar partyRegistrar.Registrar, rootPasswordLocation string) error {
	// try and retrieve the root system entity
	var systemEntityCreatedOrRetrieved system.System
	systemEntityRetrieveResponse := systemRecordHandler.RetrieveResponse{}
	err := handler.Retrieve(&systemRecordHandler.RetrieveRequest{
		Identifier: name.Identifier{Name: systemEntity.Name},
	},
		&systemEntityRetrieveResponse)
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
			defaultSystemPassword = strings.TrimSuffix(string(pwd), "\n")
		}

		// now try create the system
		systemEntityCreateResponse := systemRecordHandler.CreateResponse{}
		if err := handler.Create(&systemRecordHandler.CreateRequest{
			System: systemEntity,
		}, &systemEntityCreateResponse);
			err != nil {
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
	if err := registrar.RegisterSystemAdminUser(&partyRegistrar.RegisterSystemAdminUserRequest{
		User:     systemAdminUser,
		Password: defaultSystemPassword,
	},
		&partyRegistrar.RegisterSystemAdminUserResponse{});
		err != nil {
		switch err.(type) {
		case partyRegistrarException.AlreadyRegistered:
			// this is fine, no issues
		default:
			// something went wrong
			return systemSetupException.InitialSetup{Reasons: []string{"registration error", err.Error()}}
		}
	}

	return nil
}

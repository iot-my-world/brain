package setup

import (
	"github.com/iot-my-world/brain/pkg/party"
	partyRegistrar "github.com/iot-my-world/brain/pkg/party/registrar"
	exception3 "github.com/iot-my-world/brain/pkg/party/registrar/exception"
	system2 "github.com/iot-my-world/brain/pkg/party/system"
	"github.com/iot-my-world/brain/pkg/party/system/recordHandler"
	exception2 "github.com/iot-my-world/brain/pkg/party/system/recordHandler/exception"
	"github.com/iot-my-world/brain/pkg/party/system/setup/exception"
	"github.com/iot-my-world/brain/pkg/search/identifier/id"
	"github.com/iot-my-world/brain/pkg/search/identifier/name"
	humanUserLoginClaims "github.com/iot-my-world/brain/pkg/security/claims/login/user/human"
	humanUser "github.com/iot-my-world/brain/pkg/user/human"
	"io/ioutil"
	"os"
	"strings"
	"time"
)

var systemEntity = system2.System{
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
	handler recordHandler.RecordHandler,
	registrar partyRegistrar.Registrar,
	rootPasswordLocation string,
	systemClaims *humanUserLoginClaims.Login,
) error {
	// try and retrieve the root system entity
	var systemEntityCreatedOrRetrieved system2.System
	systemEntityRetrieveResponse, err := handler.Retrieve(&recordHandler.RetrieveRequest{
		Claims:     systemClaims,
		Identifier: name.Identifier{Name: systemEntity.Name},
	})
	switch err.(type) {
	case nil:
		// this means the system entity already exists
		systemEntityCreatedOrRetrieved = systemEntityRetrieveResponse.System

	case exception2.NotFound:
		// this means that system must be created now

		// get the password from file if a path is provided
		if rootPasswordLocation != "" {
			pwd, err := consumePasswordFile(rootPasswordLocation)
			if err != nil {
				return exception.InitialSetup{Reasons: []string{"consume password error", err.Error()}}
			}
			systemAdminUser.Password = []byte(strings.TrimSuffix(string(pwd), "\n"))
		}

		// now try create the system
		systemEntityCreateResponse, err := handler.Create(&recordHandler.CreateRequest{
			System: systemEntity,
		})
		if err != nil {
			return exception.InitialSetup{Reasons: []string{"create error", err.Error()}}
		}
		systemEntityCreatedOrRetrieved = systemEntityCreateResponse.System

	default:
		// some other error
		return exception.InitialSetup{Reasons: []string{"retrieve error", err.Error()}}
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
		case exception3.AlreadyRegistered:
			// this is fine, no issues
		default:
			// something went wrong
			return exception.InitialSetup{Reasons: []string{"registration error", err.Error()}}
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

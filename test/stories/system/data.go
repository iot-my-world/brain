package system

import (
	humanUser "github.com/iot-my-world/brain/pkg/user/human"
)

var User = humanUser.User{
	//Id string

	// Personal Details
	Name:    "root",
	Surname: "root",

	// System Details
	Username: "root",
	// EmailAddress
	Password: []byte("12345"),
	Roles:    make([]string, 0),

	// Party Details
	//ParentPartyType
	//ParentId
	//PartyType
	//PartyId
}

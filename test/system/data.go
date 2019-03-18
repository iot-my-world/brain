package system

import "gitlab.com/iotTracker/brain/party/user"

var User = user.User{
	//Id string

	// Personal Details
	Name:    "root",
	Surname: "root",

	// System Details
	Username: "root",
	// EmailAddress
	Password: []byte("12345"),
	// Roles

	// Party Details
	//ParentPartyType
	//ParentId
	//PartyType
	//PartyId
}

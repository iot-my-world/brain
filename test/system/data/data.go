package data

import "gitlab.com/iotTracker/brain/user"

var User = user.User{
	//Id string

	// Personal Details
	Name:    "root",
	Surname: "root",

	// System Details
	Username: "root",
	// EmailAddress
	Password: []byte("12345"),
	//Password: []byte("thebrainstemistherootofallthought"), // for server testing
	// Roles

	// Party Details
	//ParentPartyType
	//ParentId
	//PartyType
	//PartyId
}
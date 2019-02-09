package user

import (
	"gitlab.com/iotTracker/brain/party"
	"gitlab.com/iotTracker/brain/search/identifier/id"
)

// Defines the User record for the database
type User struct {
	Id string `json:"id" bson:"id"`

	// Personal Details
	Name    string `json:"name" bson:"name"`
	Surname string `json:"surname" bson:"surname"`

	// System Details
	Username     string   `json:"username" bson:"username"`
	EmailAddress string   `json:"emailAddress" bson:"emailAddress"`
	Password     []byte   `json:"pwd" bson:"pwd"`
	Roles        []string `json:"roles" bson:"roles"`

	// Party Details
	PartyType party.Type    `json:"partyType" bson:"partyType"`
	PartyId   id.Identifier `json:"partyId" bson:"partyId"`
}

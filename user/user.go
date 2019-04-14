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
	Password     []byte   `json:"password" bson:"password"`
	Roles        []string `json:"roles" bson:"roles"`

	// Party Details
	ParentPartyType party.Type    `json:"parentPartyType" bson:"parentPartyType"`
	ParentId        id.Identifier `json:"parentId" bson:"parentId"`
	PartyType       party.Type    `json:"partyType" bson:"partyType"`
	PartyId         id.Identifier `json:"partyId" bson:"partyId"`

	Registered bool `json:"registered" bson:"registered"`
}

func SetId() {

}

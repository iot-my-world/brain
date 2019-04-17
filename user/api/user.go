package api

import (
	"gitlab.com/iotTracker/brain/party"
	"gitlab.com/iotTracker/brain/search/identifier/id"
)

type User struct {
	Id string `json:"id" bson:"id"`

	// Details
	Name        string `json:"name" bson:"name"`
	Description string `json:"description" bson:"description"`

	// System Details
	Username string   `json:"username" bson:"username"`
	Password []byte   `json:"password" bson:"password"`
	Roles    []string `json:"roles" bson:"roles"`

	// Party Details
	PartyType party.Type    `json:"parentPartyType" bson:"parentPartyType"`
	PartyId   id.Identifier `json:"partyId" bson:"partyId"`
}

func (u *User) SetId(id string) {
	u.Id = id
}

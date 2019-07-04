package api

import (
	"github.com/iot-my-world/brain/pkg/party"
	"github.com/iot-my-world/brain/pkg/search/identifier/id"
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
	PartyType party.Type    `json:"partyType" bson:"partyType"`
	PartyId   id.Identifier `json:"partyId" bson:"partyId"`
}

func (u *User) SetId(id string) {
	u.Id = id
}

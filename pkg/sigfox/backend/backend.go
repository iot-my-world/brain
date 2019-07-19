package backend

import (
	"github.com/iot-my-world/brain/pkg/party"
	"github.com/iot-my-world/brain/pkg/search/identifier/id"
)

type Backend struct {
	Id string `json:"id" bson:"id"`

	OwnerPartyType party.Type    `json:"ownerPartyType" bson:"ownerPartyType"`
	OwnerId        id.Identifier `json:"ownerId" bson:"ownerId"`
	Name           string        `json:"name" bson:"name"`
	Token          string        `json:"token" bson:"token"`
}

func (b *Backend) SetId(id string) {
	b.Id = id
}

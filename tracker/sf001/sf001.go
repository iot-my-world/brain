package sf001

import (
	"github.com/iot-my-world/brain/party"
	"github.com/iot-my-world/brain/search/identifier/id"
)

type SF001 struct {
	Id       string `json:"id" bson:"id"`
	DeviceId string `json:"deviceId" bson:"deviceId"`

	OwnerPartyType    party.Type    `json:"ownerPartyType" bson:"ownerPartyType"`
	OwnerId           id.Identifier `json:"ownerId" bson:"ownerId"`
	AssignedPartyType party.Type    `json:"assignedPartyType" bson:"assignedPartyType"`
	AssignedId        id.Identifier `json:"assignedId" bson:"assignedId"`

	LastMessageTimestamp int64 `json:"lastMessageTimestamp" bson:"lastMessageTimestamp"`
}

func (s *SF001) SetId(id string) {
	s.Id = id
}

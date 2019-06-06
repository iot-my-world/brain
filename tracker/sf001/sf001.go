package sf001

import (
	"gitlab.com/iotTracker/brain/party"
	"gitlab.com/iotTracker/brain/search/identifier/id"
)

type SF001 struct {
	Id       string `json:"string" bson:"string"`
	DeviceId string `json:"string" bson:"string"`

	OwnerPartyType    party.Type    `json:"ownerPartyType" bson:"ownerPartyType"`
	OwnerId           id.Identifier `json:"ownerId" bson:"ownerId"`
	AssignedPartyType party.Type    `json:"assignedPartyType" bson:"assignedPartyType"`
	AssignedId        id.Identifier `json:"assignedId" bson:"assignedId"`

	LastMessageTimestamp int64 `json:"lastMessageTimestamp" bson:"lastMessageTimestamp"`
}

func (s *SF001) SetId(id string) {
	s.Id = id
}

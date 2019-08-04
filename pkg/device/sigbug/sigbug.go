package sigbug

import (
	"github.com/iot-my-world/brain/pkg/party"
	"github.com/iot-my-world/brain/pkg/search/identifier/id"
	sigfoxBackendDataCallbackMessage "github.com/iot-my-world/brain/pkg/sigfox/backend/callback/data/message"
)

type Sigbug struct {
	Id       string `json:"id" bson:"id"`
	DeviceId string `json:"deviceId" bson:"deviceId"`

	OwnerPartyType    party.Type    `json:"ownerPartyType" bson:"ownerPartyType"`
	OwnerId           id.Identifier `json:"ownerId" bson:"ownerId"`
	AssignedPartyType party.Type    `json:"assignedPartyType" bson:"assignedPartyType"`
	AssignedId        id.Identifier `json:"assignedId" bson:"assignedId"`

	LastMessage sigfoxBackendDataCallbackMessage.Message `json:"lastMessage" bson:"lastMessage"`
}

func (s *Sigbug) SetId(id string) {
	s.Id = id
}

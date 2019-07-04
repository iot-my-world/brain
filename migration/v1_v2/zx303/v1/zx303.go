package v1

import (
	"github.com/iot-my-world/brain/pkg/party"
	"github.com/iot-my-world/brain/search/identifier/id"
)

type ZX303 struct {
	Id                string        `json:"id" bson:"id"`
	IMEI              string        `json:"imei" bson:"imei"`
	SimCountryCode    string        `json:"simCountryCode" bson:"simCountryCode"`
	SimNumber         string        `json:"simNumber" bson:"simNumber"`
	OwnerPartyType    party.Type    `json:"ownerPartyType" bson:"ownerPartyType"`
	OwnerId           id.Identifier `json:"ownerId" bson:"ownerId"`
	AssignedPartyType party.Type    `json:"assignedPartyType" bson:"assignedPartyType"`
	AssignedId        id.Identifier `json:"assignedId" bson:"assignedId"`
}

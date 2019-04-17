package tk102

import (
	"gitlab.com/iotTracker/brain/party"
	"gitlab.com/iotTracker/brain/search/identifier/id"
	"gitlab.com/iotTracker/brain/tracker/device"
)

type TK102 struct {
	Type              device.Type   `json:"type" bson:"type"`
	Id                string        `json:"id" bson:"id"`
	ManufacturerId    string        `json:"manufacturerId" bson:"manufacturerId"`
	SimCountryCode    string        `json:"simCountryCode" bson:"simCountryCode"`
	SimNumber         string        `json:"simNumber" bson:"simNumber"`
	OwnerPartyType    party.Type    `json:"ownerPartyType" bson:"ownerPartyType"`
	OwnerId           id.Identifier `json:"ownerId" bson:"ownerId"`
	AssignedPartyType party.Type    `json:"assignedPartyType" bson:"assignedPartyType"`
	AssignedId        id.Identifier `json:"assignedId" bson:"assignedId"`
}

func (t TK102) DeviceType() device.Type {
	return t.Type
}

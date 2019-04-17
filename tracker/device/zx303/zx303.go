package zx303

import (
	"gitlab.com/iotTracker/brain/party"
	"gitlab.com/iotTracker/brain/search/identifier/id"
	"gitlab.com/iotTracker/brain/tracker/device"
)

type ZX303 struct {
	Type              device.Type   `json:"type" bson:"type"`
	Id                string        `json:"id" bson:"id"`
	IMEI              string        `json:"imei" bson:"imei"`
	SimCountryCode    string        `json:"simCountryCode" bson:"simCountryCode"`
	SimNumber         string        `json:"simNumber" bson:"simNumber"`
	OwnerPartyType    party.Type    `json:"ownerPartyType" bson:"ownerPartyType"`
	OwnerId           id.Identifier `json:"ownerId" bson:"ownerId"`
	AssignedPartyType party.Type    `json:"assignedPartyType" bson:"assignedPartyType"`
	AssignedId        id.Identifier `json:"assignedId" bson:"assignedId"`
}

func (z *ZX303) SetId(id string) {
	z.Id = id
}

func (z ZX303) DeviceType() device.Type {
	return z.Type
}

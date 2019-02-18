package device

import (
	"gitlab.com/iotTracker/brain/party"
	"gitlab.com/iotTracker/brain/search/identifier/id"
)

type Device struct {
	Id                string        `json:"id" bson:"id"`
	IMEI              string        `json:"imei" bson:"imei"`
	SimCountryCode    string        `json:"simCountryCode" bson:"simCountryCode"`
	SimNumber         string        `json:"simNumber" bson:"simNumber"`
	OwnerPartyType    party.Type    `json:"ownerPartyType" bson:"ownerPartyType"`
	OwnerId           id.Identifier `json:"ownerId" bson:"ownerId"`
	AssignedPartyType party.Type    `json:"assignedPartyType" bson:"assignedPartyType"`
	AssignedId        id.Identifier `json:"assignedId" bson:"assignedId"`
}

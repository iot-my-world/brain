package v2

import (
	"github.com/iot-my-world/brain/pkg/party"
	"github.com/iot-my-world/brain/pkg/search/identifier/id"
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

	LoggedIn               bool  `json:"loggedIn" bson:"loggedIn"`
	LogInTimestamp         int64 `json:"logInTimestamp" bson:"logInTimestamp"`
	LogOutTimestamp        int64 `json:"logOutTimestamp" bson:"logOutTimestamp"`
	LastHeartbeatTimestamp int64 `json:"lastHeartbeatTimestamp" bson:"lastHeartbeatTimestamp"`
}

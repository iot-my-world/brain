package zx303

import (
	"encoding/json"
	"errors"
	"gitlab.com/iotTracker/brain/party"
	"gitlab.com/iotTracker/brain/search/identifier"
	wrappedIdentifier "gitlab.com/iotTracker/brain/search/identifier/wrapped"
	"gitlab.com/iotTracker/brain/tracker/device"
)

type ZX303 struct {
	Id                string                `json:"id" bson:"id"`
	IMEI              string                `json:"imei"`
	SimCountryCode    string                `json:"simCountryCode" bson:"simCountryCode"`
	SimNumber         string                `json:"simNumber" bson:"simNumber"`
	OwnerPartyType    party.Type            `json:"ownerPartyType" bson:"ownerPartyType"`
	OwnerId           identifier.Identifier `json:"ownerId" bson:"ownerId"`
	AssignedPartyType party.Type            `json:"assignedPartyType" bson:"assignedPartyType"`
	AssignedId        identifier.Identifier `json:"assignedId" bson:"assignedId"`
}

func (z ZX303) Type() device.Type {
	return device.ZX303
}

func (z *ZX303) UnmarshalJSON(data []byte) error {
	type Alias ZX303
	aux := &struct {
		OwnerId    wrappedIdentifier.Wrapped `json:"ownerId"`
		AssignedId wrappedIdentifier.Wrapped `json:"assignedId"`
		*Alias
	}{
		Alias: (*Alias)(z),
	}

	if err := json.Unmarshal(data, &aux); err != nil {
		return errors.New("unwrapping: " + err.Error())
	}
	z.OwnerId = aux.OwnerId.Identifier
	z.AssignedId = aux.AssignedId.Identifier

	return nil
}

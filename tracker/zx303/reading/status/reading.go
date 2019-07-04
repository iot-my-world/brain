package status

import (
	"github.com/iot-my-world/brain/pkg/party"
	"github.com/iot-my-world/brain/search/identifier"
	"github.com/iot-my-world/brain/search/identifier/id"
)

type Reading struct {
	Id string `json:"id" bson:"id"`

	// Device Details
	DeviceId id.Identifier `json:"deviceId" bson:"deviceId"`

	// Owner Details
	OwnerPartyType    party.Type    `json:"ownerPartyType" bson:"ownerPartyType"`
	OwnerId           id.Identifier `json:"ownerId" bson:"ownerId"`
	AssignedPartyType party.Type    `json:"assignedPartyType" bson:"assignedPartyType"`
	AssignedId        id.Identifier `json:"assignedId" bson:"assignedId"`

	// Reading Details
	Timestamp         int64 `json:"timestamp" bson:"timeStamp"`
	BatteryPercentage int64 `json:"batteryPercentage" bson:"batteryPercentage"`
	UploadInterval    int64 `json:"uploadInterval" bson:"uploadInterval"`
	SoftwareVersion   int64 `json:"softwareVersion" bson:"softwareVersion"`
	Timezone          int64 `json:"timezone" bson:"timezone"`
}

func (r *Reading) SetId(id string) {
	r.Id = id
}

func IsValidIdentifier(id identifier.Identifier) bool {
	if id == nil {
		return false
	}

	switch id.Type() {
	case identifier.Id:
		return true
	default:
		return false
	}
}

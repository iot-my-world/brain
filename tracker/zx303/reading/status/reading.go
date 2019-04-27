package reading

import (
	"gitlab.com/iotTracker/brain/search/identifier"
)

type Reading struct {
	Id string `json:"id" bson:"id"`

	Timestamp         int64 `json:"timestamp" bson:"timeStamp"`
	BatteryPercentage int64 `json:"batteryPercentage" bson:"batteryPercentage"`
	UploadInterval    int64 `json:"uploadInterval" bson:"uploadInterval"`
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

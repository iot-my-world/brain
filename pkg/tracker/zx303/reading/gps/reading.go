package gps

import (
	"github.com/iot-my-world/brain/pkg/party"
	"github.com/iot-my-world/brain/pkg/search/identifier"
	"github.com/iot-my-world/brain/pkg/search/identifier/id"
	"math"
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
	TimeStamp    int64   `json:"timeStamp" bson:"timeStamp"`
	NoSatellites int64   `json:"noSatellites" bson:"noSatellites"`
	Latitude     float32 `json:"latitude" bson:"latitude"`
	Longitude    float32 `json:"longitude" bson:"longitude"`
	Speed        int64   `json:"speed" bson:"speed"`
	Heading      int64   `json:"heading" bson:"heading"`
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

const earthRadiusInKm float64 = 6378.137

func DifferenceBetween(r1, r2 *Reading) float32 {
	lat1 := r1.Latitude
	lon1 := r1.Longitude
	lat2 := r2.Latitude
	lon2 := r2.Longitude
	var dLat = float64(lat2*math.Pi/180 - lat1*math.Pi/180)
	var dLon = float64(lon2*math.Pi/180 - lon1*math.Pi/180)
	var a = math.Sin(dLat/2)*math.Sin(dLat/2) +
		math.Cos(float64(lat1*math.Pi/180))*math.Cos(float64(lat2*math.Pi/180))*
			math.Sin(dLon/2)*math.Sin(dLon/2)
	var c = 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))
	var d = earthRadiusInKm * c
	return float32(d * 1000)
}

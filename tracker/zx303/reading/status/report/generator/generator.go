package generator

import (
	"gitlab.com/iotTracker/brain/search/identifier"
	"gitlab.com/iotTracker/brain/security/claims"
)

type Generator interface {
	Battery(request *BatteryRequest) (*BatteryResponse, error)
}

type BatteryRequest struct {
	Claims                 claims.Claims
	ZX303TrackerIdentifier identifier.Identifier
}

type BatteryResponse struct {
	Readings []BatteryReading
}

type BatteryReading struct {
	Percentage int64 `json:"batteryPercentage"`
	Timestamp  int64 `json:"timestamp"`
}

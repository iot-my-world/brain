package tracking

import (
	"github.com/iot-my-world/brain/search/identifier/party"
	"github.com/iot-my-world/brain/security/claims"
	zx303TrackerGPSReading "github.com/iot-my-world/brain/tracker/zx303/reading/gps"
)

type Report interface {
	Live(request *LiveRequest) (*LiveResponse, error)
	Historical(request *HistoricalRequest) (*HistoricalResponse, error)
}

type LiveRequest struct {
	Claims           claims.Claims
	PartyIdentifiers []party.Identifier
}

type LiveResponse struct {
	ZX303TrackerGPSReadings []zx303TrackerGPSReading.Reading
}

type HistoricalRequest struct {
	Claims claims.Claims
}

type HistoricalResponse struct {
	ZX303TrackerGPSReadings []zx303TrackerGPSReading.Reading
}

package tracking

import (
	"gitlab.com/iotTracker/brain/search/identifier/party"
	"gitlab.com/iotTracker/brain/security/claims"
	zx303TrackerGPSReading "gitlab.com/iotTracker/brain/tracker/zx303/reading/gps"
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

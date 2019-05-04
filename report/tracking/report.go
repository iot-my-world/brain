package tracking

import (
	"gitlab.com/iotTracker/brain/search/identifier/party"
	"gitlab.com/iotTracker/brain/security/claims"
	"gitlab.com/iotTracker/brain/tracker/tk102/reading"
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
	Readings []reading.Reading
}

type HistoricalRequest struct {
	Claims claims.Claims
}

type HistoricalResponse struct {
	Readings []reading.Reading
}

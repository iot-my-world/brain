package tracking

import (
	"gitlab.com/iotTracker/brain/search/identifier"
	"gitlab.com/iotTracker/brain/security/claims"
	"gitlab.com/iotTracker/brain/tracker/reading"
)

type Report interface {
	Live(request *LiveRequest, response *LiveResponse) error
	Historical(request *HistoricalRequest, response *HistoricalResponse) error
}

type LiveRequest struct {
	Claims             claims.Claims
	ClientIdentifiers  []identifier.Identifier
	CompanyIdentifiers []identifier.Identifier
}

type LiveResponse struct {
	Readings []reading.Reading
}

type HistoricalRequest struct {
	Claims             claims.Claims
	ClientIdentifiers  []identifier.Identifier
	CompanyIdentifiers []identifier.Identifier
}

type HistoricalResponse struct {
	Readings []reading.Reading
}

package jsonRpc

import (
	"gitlab.com/iotTracker/brain/log"
	trackingReport "gitlab.com/iotTracker/brain/report/tracking"
	"gitlab.com/iotTracker/brain/search/identifier"
	"gitlab.com/iotTracker/brain/search/identifier/party"
	wrappedIdentifier "gitlab.com/iotTracker/brain/search/identifier/wrapped"
	"gitlab.com/iotTracker/brain/security/wrappedClaims"
	"gitlab.com/iotTracker/brain/tracker/reading"
	"net/http"
)

type adaptor struct {
	trackingReport trackingReport.Report
}

func New(
	trackingReport trackingReport.Report,
) *adaptor {
	return &adaptor{
		trackingReport: trackingReport,
	}
}

type LiveRequest struct {
	PartyIdentifiers []party.Identifier `json:"partyIdentifiers"`
}

type LiveResponse struct {
	Readings []reading.Reading `json:"readings"`
}

func (a *adaptor) Live(r *http.Request, request *LiveRequest, response *LiveResponse) error {
	claims, err := wrappedClaims.UnwrapClaimsFromContext(r)
	if err != nil {
		log.Warn(err.Error())
		return err
	}

	// get report
	liveTrackingReportResponse, err := a.trackingReport.Live(&trackingReport.LiveRequest{
		Claims:           claims,
		PartyIdentifiers: request.PartyIdentifiers,
	})
	if err != nil {
		return err
	}

	response.Readings = liveTrackingReportResponse.Readings

	return nil
}

type HistoricalRequest struct {
	CompanyIdentifiers []wrappedIdentifier.Wrapped `json:"companyIdentifiers"`
	ClientIdentifiers  []wrappedIdentifier.Wrapped `json:"clientIdentifiers"`
}

type HistoricalResponse struct {
	Readings []reading.Reading `json:"readings"`
}

func (a *adaptor) Historical(r *http.Request, request *HistoricalRequest, response *HistoricalResponse) error {
	claims, err := wrappedClaims.UnwrapClaimsFromContext(r)
	if err != nil {
		log.Warn(err.Error())
		return err
	}

	// unwrap company identifiers
	companyIdentifiers := make([]identifier.Identifier, 0)
	for idIdx := range request.CompanyIdentifiers {
		if c, err := request.CompanyIdentifiers[idIdx].UnWrap(); err == nil {
			companyIdentifiers = append(companyIdentifiers, c)
		} else {
			return err
		}
	}
	// unwrap client criteria
	clientIdentifiers := make([]identifier.Identifier, 0)
	for idIdx := range request.ClientIdentifiers {
		if c, err := request.ClientIdentifiers[idIdx].UnWrap(); err == nil {
			clientIdentifiers = append(clientIdentifiers, c)
		} else {
			return err
		}
	}

	// get report
	historicalTrackingReportResponse, err := a.trackingReport.Historical(&trackingReport.HistoricalRequest{
		Claims: claims,
	})
	if err != nil {
		return err
	}

	response.Readings = historicalTrackingReportResponse.Readings

	return nil
}

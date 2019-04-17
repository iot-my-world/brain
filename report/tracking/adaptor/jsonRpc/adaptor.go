package jsonRpc

import (
	"errors"
	"gitlab.com/iotTracker/brain/log"
	trackingReport "gitlab.com/iotTracker/brain/report/tracking"
	"gitlab.com/iotTracker/brain/search/identifier"
	"gitlab.com/iotTracker/brain/search/identifier/party"
	wrappedIdentifier "gitlab.com/iotTracker/brain/search/identifier/wrapped"
	wrappedClaims "gitlab.com/iotTracker/brain/security/claims/wrapped"
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
	WrappedPartyIdentifiers []wrappedIdentifier.Wrapped `json:"partyIdentifiers"`
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

	partyIdentifiers := make([]party.Identifier, 0)

	for i := range request.WrappedPartyIdentifiers {
		partyIdentifier, ok := request.WrappedPartyIdentifiers[i].Identifier.(party.Identifier)
		if !ok {
			return errors.New("could not cast identifier.Identifier to party.Identifier")
		}
		partyIdentifiers = append(partyIdentifiers, partyIdentifier)
	}

	// get report
	liveTrackingReportResponse, err := a.trackingReport.Live(&trackingReport.LiveRequest{
		Claims:           claims,
		PartyIdentifiers: partyIdentifiers,
	})
	if err != nil {
		return err
	}

	response.Readings = liveTrackingReportResponse.Readings

	return nil
}

type HistoricalRequest struct {
	WrappedCompanyIdentifiers []wrappedIdentifier.Wrapped `json:"companyIdentifiers"`
	WrappedClientIdentifiers  []wrappedIdentifier.Wrapped `json:"clientIdentifiers"`
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
	for i := range request.WrappedCompanyIdentifiers {
		partyIdentifier, ok := request.WrappedCompanyIdentifiers[i].Identifier.(party.Identifier)
		if !ok {
			return errors.New("could not cast identifier.Identifier to party.Identifier")
		}
		companyIdentifiers = append(companyIdentifiers, partyIdentifier)
	}
	// unwrap client criteria
	clientIdentifiers := make([]identifier.Identifier, 0)
	for i := range request.WrappedClientIdentifiers {
		partyIdentifier, ok := request.WrappedClientIdentifiers[i].Identifier.(party.Identifier)
		if !ok {
			return errors.New("could not cast identifier.Identifier to party.Identifier")
		}
		clientIdentifiers = append(clientIdentifiers, partyIdentifier)
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

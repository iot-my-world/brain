package jsonRpc

import (
	"errors"
	"github.com/iot-my-world/brain/log"
	trackingReport "github.com/iot-my-world/brain/report/tracking"
	"github.com/iot-my-world/brain/search/identifier"
	"github.com/iot-my-world/brain/search/identifier/party"
	wrappedIdentifier "github.com/iot-my-world/brain/search/identifier/wrapped"
	wrappedClaims "github.com/iot-my-world/brain/security/claims/wrapped"
	zx303TrackerGPSReading "github.com/iot-my-world/brain/tracker/zx303/reading/gps"
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
	ZX303TrackerGPSReadings []zx303TrackerGPSReading.Reading `json:"zx303TrackerGPSReadings"`
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

	response.ZX303TrackerGPSReadings = liveTrackingReportResponse.ZX303TrackerGPSReadings

	return nil
}

type HistoricalRequest struct {
	WrappedCompanyIdentifiers []wrappedIdentifier.Wrapped `json:"companyIdentifiers"`
	WrappedClientIdentifiers  []wrappedIdentifier.Wrapped `json:"clientIdentifiers"`
}

type HistoricalResponse struct {
	ZX303TrackerGPSReadings []zx303TrackerGPSReading.Reading `json:"zx303TrackerGPSReadings"`
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

	response.ZX303TrackerGPSReadings = historicalTrackingReportResponse.ZX303TrackerGPSReadings

	return nil
}

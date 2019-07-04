package jsonRpc

import (
	"errors"
	"github.com/iot-my-world/brain/internal/log"
	"github.com/iot-my-world/brain/pkg/report/tracking"
	"github.com/iot-my-world/brain/pkg/search/identifier"
	"github.com/iot-my-world/brain/pkg/search/identifier/party"
	wrappedIdentifier "github.com/iot-my-world/brain/pkg/search/identifier/wrapped"
	wrappedClaims "github.com/iot-my-world/brain/pkg/security/claims/wrapped"
	zx303TrackerGPSReading "github.com/iot-my-world/brain/pkg/tracker/zx303/reading/gps"
	"net/http"
)

type adaptor struct {
	trackingReport tracking.Report
}

func New(
	trackingReport tracking.Report,
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
	liveTrackingReportResponse, err := a.trackingReport.Live(&tracking.LiveRequest{
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
	historicalTrackingReportResponse, err := a.trackingReport.Historical(&tracking.HistoricalRequest{
		Claims: claims,
	})
	if err != nil {
		return err
	}

	response.ZX303TrackerGPSReadings = historicalTrackingReportResponse.ZX303TrackerGPSReadings

	return nil
}

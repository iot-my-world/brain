package jsonRpc

import (
	trackingReport "gitlab.com/iotTracker/brain/report/tracking"
	"gitlab.com/iotTracker/brain/tracker/reading"
	"net/http"
	"gitlab.com/iotTracker/brain/security/wrappedClaims"
	"gitlab.com/iotTracker/brain/log"
	"gitlab.com/iotTracker/brain/search/wrappedIdentifier"
	"gitlab.com/iotTracker/brain/search/identifier"
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
	CompanyIdentifiers []wrappedIdentifier.WrappedIdentifier `json:"companyIdentifiers"`
	ClientIdentifiers  []wrappedIdentifier.WrappedIdentifier `json:"clientIdentifiers"`
}

type LiveResponse struct {
	Readings []reading.Reading `json:"readings"`
}

func (a *adaptor) Live(r *http.Request, request *LiveRequest, response *LiveResponse) error {
	// unwrap claims
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
	liveTrackingReportResponse := trackingReport.LiveResponse{}
	if err := a.trackingReport.Live(&trackingReport.LiveRequest{
		Claims:             claims,
		CompanyIdentifiers: companyIdentifiers,
		ClientIdentifiers:  clientIdentifiers,
	}, &liveTrackingReportResponse); err != nil {
		return err
	}

	response.Readings = liveTrackingReportResponse.Readings

	return nil
}

type HistoricalRequest struct {
	CompanyIdentifiers []wrappedIdentifier.WrappedIdentifier `json:"companyIdentifiers"`
	ClientIdentifiers  []wrappedIdentifier.WrappedIdentifier `json:"clientIdentifiers"`
}

type HistoricalResponse struct {
	Readings []reading.Reading `json:"readings"`
}

func (a *adaptor) Historical(r *http.Request, request *HistoricalRequest, response *HistoricalResponse) error {
	// unwrap claims
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
	liveTrackingReportResponse := trackingReport.LiveResponse{}
	if err := a.trackingReport.Live(&trackingReport.LiveRequest{
		Claims:             claims,
		CompanyIdentifiers: companyIdentifiers,
		ClientIdentifiers:  clientIdentifiers,
	}, &liveTrackingReportResponse); err != nil {
		return err
	}

	return nil
}

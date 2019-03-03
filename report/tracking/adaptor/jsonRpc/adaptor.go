package jsonRpc

import (
	trackingReport "gitlab.com/iotTracker/brain/report/tracking"
	"gitlab.com/iotTracker/brain/party/company"
	"gitlab.com/iotTracker/brain/tracker/reading"
	"gitlab.com/iotTracker/brain/search/wrappedCriterion"
	"gitlab.com/iotTracker/brain/party/client"
	"net/http"
	"gitlab.com/iotTracker/brain/search/criterion"
	"gitlab.com/iotTracker/brain/security/wrappedClaims"
	"gitlab.com/iotTracker/brain/log"
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
	CompanyCriteria []wrappedCriterion.WrappedCriterion `json:"companyCriteria"`
	ClientCriteria  []wrappedCriterion.WrappedCriterion `json:"clientCriteria"`
}

type LiveResponse struct {
	Companies []company.Company `json:"companies"`
	Clients   []client.Client   `json:"clients"`
	Readings  []reading.Reading `json:"readings"`
}

func (a *adaptor) Live(r *http.Request, request *LiveRequest, response *LiveResponse) error {
	// unwrap claims
	claims, err := wrappedClaims.UnwrapClaimsFromContext(r)
	if err != nil {
		log.Warn(err.Error())
		return err
	}

	// unwrap company criteria
	companyCriteria := make([]criterion.Criterion, 0)
	for criterionIdx := range request.CompanyCriteria {
		if c, err := request.CompanyCriteria[criterionIdx].UnWrap(); err == nil {
			companyCriteria = append(companyCriteria, c)
		} else {
			return err
		}
	}
	// unwrap client criteria
	clientCriteria := make([]criterion.Criterion, 0)
	for criterionIdx := range request.ClientCriteria {
		if c, err := request.ClientCriteria[criterionIdx].UnWrap(); err == nil {
			clientCriteria = append(clientCriteria, c)
		} else {
			return err
		}
	}

	// get report
	liveTrackingReportResponse := trackingReport.LiveResponse{}
	if err := a.trackingReport.Live(&trackingReport.LiveRequest{
		Claims:          claims,
		ClientCriteria:  clientCriteria,
		CompanyCriteria: companyCriteria,
	}, &liveTrackingReportResponse); err != nil {
		return err
	}

	return nil
}

type HistoricalRequest struct {
	CompanyCriteria []wrappedCriterion.WrappedCriterion `json:"companyCriteria"`
	ClientCriteria  []wrappedCriterion.WrappedCriterion `json:"clientCriteria"`
}

type HistoricalResponse struct {
	Companies []company.Company `json:"companies"`
	Clients   []client.Client   `json:"clients"`
	Readings  []reading.Reading `json:"readings"`
}

func (a *adaptor) Historical(r *http.Request, request *HistoricalRequest, response *HistoricalResponse) error {
	// unwrap claims
	claims, err := wrappedClaims.UnwrapClaimsFromContext(r)
	if err != nil {
		log.Warn(err.Error())
		return err
	}

	// unwrap company criteria
	companyCriteria := make([]criterion.Criterion, 0)
	for criterionIdx := range request.CompanyCriteria {
		if c, err := request.CompanyCriteria[criterionIdx].UnWrap(); err == nil {
			companyCriteria = append(companyCriteria, c)
		} else {
			return err
		}
	}
	// unwrap client criteria
	clientCriteria := make([]criterion.Criterion, 0)
	for criterionIdx := range request.ClientCriteria {
		if c, err := request.ClientCriteria[criterionIdx].UnWrap(); err == nil {
			clientCriteria = append(clientCriteria, c)
		} else {
			return err
		}
	}

	// get report
	historicalTrackingReportResponse := trackingReport.HistoricalResponse{}
	if err := a.trackingReport.Historical(&trackingReport.HistoricalRequest{
		Claims:          claims,
		ClientCriteria:  clientCriteria,
		CompanyCriteria: companyCriteria,
	}, &historicalTrackingReportResponse); err != nil {
		return err
	}

	return nil
}

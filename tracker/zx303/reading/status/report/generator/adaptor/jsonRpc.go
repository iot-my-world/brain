package adaptor

import (
	"gitlab.com/iotTracker/brain/log"
	wrappedIdentifier "gitlab.com/iotTracker/brain/search/identifier/wrapped"
	wrappedClaims "gitlab.com/iotTracker/brain/security/claims/wrapped"
	zx303StatusReadingReport "gitlab.com/iotTracker/brain/tracker/zx303/reading/status/report"
	zx303StatusReadingReportGenerator "gitlab.com/iotTracker/brain/tracker/zx303/reading/status/report/generator"
	"net/http"
)

type Adaptor struct {
	zx303StatusReadingReportGenerator zx303StatusReadingReportGenerator.Generator
}

func New(
	zx303StatusReadingReportGenerator zx303StatusReadingReportGenerator.Generator,
) *Adaptor {
	return &Adaptor{
		zx303StatusReadingReportGenerator: zx303StatusReadingReportGenerator,
	}
}

type BatteryReportRequest struct {
	ZX303TrackerIdentifier wrappedIdentifier.Wrapped `json:"zx303TrackerIdentifier"`
	StartDate              int64                     `json:"startDate"`
	EndDate                int64                     `json:"endDate"`
}

type BatteryReportResponse struct {
	Report zx303StatusReadingReport.Battery `json:"report"`
}

func (a *Adaptor) BatteryReport(r *http.Request, request *BatteryReportRequest, response *BatteryReportResponse) error {
	claims, err := wrappedClaims.UnwrapClaimsFromContext(r)
	if err != nil {
		log.Warn(err.Error())
		return err
	}

	batteryStatusReportResponse, err := a.zx303StatusReadingReportGenerator.BatteryReport(&zx303StatusReadingReportGenerator.BatteryReportRequest{
		Claims:                 claims,
		ZX303TrackerIdentifier: request.ZX303TrackerIdentifier.Identifier,
		StartDate:              request.StartDate,
		EndDate:                request.EndDate,
	})
	if err != nil {
		return err
	}

	response.Report = batteryStatusReportResponse.Report

	return nil
}

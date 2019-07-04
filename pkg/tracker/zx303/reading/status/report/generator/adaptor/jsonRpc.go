package adaptor

import (
	"github.com/iot-my-world/brain/internal/log"
	wrappedIdentifier "github.com/iot-my-world/brain/pkg/search/identifier/wrapped"
	wrappedClaims "github.com/iot-my-world/brain/pkg/security/claims/wrapped"
	"github.com/iot-my-world/brain/pkg/tracker/zx303/reading/status/report"
	"github.com/iot-my-world/brain/pkg/tracker/zx303/reading/status/report/generator"
	"net/http"
)

type Adaptor struct {
	zx303StatusReadingReportGenerator generator.Generator
}

func New(
	zx303StatusReadingReportGenerator generator.Generator,
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
	Report report.Battery `json:"report"`
}

func (a *Adaptor) BatteryReport(r *http.Request, request *BatteryReportRequest, response *BatteryReportResponse) error {
	claims, err := wrappedClaims.UnwrapClaimsFromContext(r)
	if err != nil {
		log.Warn(err.Error())
		return err
	}

	batteryStatusReportResponse, err := a.zx303StatusReadingReportGenerator.BatteryReport(&generator.BatteryReportRequest{
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
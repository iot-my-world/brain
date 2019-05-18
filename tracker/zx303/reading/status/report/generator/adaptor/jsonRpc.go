package adaptor

import (
	"gitlab.com/iotTracker/brain/log"
	wrappedIdentifier "gitlab.com/iotTracker/brain/search/identifier/wrapped"
	wrappedClaims "gitlab.com/iotTracker/brain/security/claims/wrapped"
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

type BatteryRequest struct {
	ZX303TrackerIdentifier wrappedIdentifier.Wrapped `json:"zx303TrackerIdentifier"`
}

type BatteryResponse struct {
	Readings []zx303StatusReadingReportGenerator.BatteryReading `json:"readings"`
}

func (a *Adaptor) Battery(r *http.Request, request *BatteryRequest, response *BatteryResponse) error {
	claims, err := wrappedClaims.UnwrapClaimsFromContext(r)
	if err != nil {
		log.Warn(err.Error())
		return err
	}

	batteryStatusReportResponse, err := a.zx303StatusReadingReportGenerator.Battery(&zx303StatusReadingReportGenerator.BatteryRequest{
		Claims:                 claims,
		ZX303TrackerIdentifier: request.ZX303TrackerIdentifier.Identifier,
	})
	if err != nil {
		return err
	}

	response.Readings = batteryStatusReportResponse.Readings

	return nil
}

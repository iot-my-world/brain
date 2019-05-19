package generator

import (
	"gitlab.com/iotTracker/brain/search/identifier"
	"gitlab.com/iotTracker/brain/security/claims"
	zx303TrackerStatusReadingReport "gitlab.com/iotTracker/brain/tracker/zx303/reading/status/report"
)

type Generator interface {
	BatteryReport(request *BatteryReportRequest) (*BatteryReportResponse, error)
}

type BatteryReportRequest struct {
	Claims                 claims.Claims
	ZX303TrackerIdentifier identifier.Identifier
	StartDate              int64
	EndDate                int64
}

type BatteryReportResponse struct {
	Report zx303TrackerStatusReadingReport.Battery
}

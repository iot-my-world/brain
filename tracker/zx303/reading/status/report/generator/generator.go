package generator

import (
	"github.com/iot-my-world/brain/search/identifier"
	"github.com/iot-my-world/brain/security/claims"
	zx303TrackerStatusReadingReport "github.com/iot-my-world/brain/tracker/zx303/reading/status/report"
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

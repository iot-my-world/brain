package generator

import (
	"github.com/iot-my-world/brain/pkg/search/identifier"
	"github.com/iot-my-world/brain/pkg/security/claims"
	"github.com/iot-my-world/brain/pkg/tracker/zx303/reading/status/report"
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
	Report report.Battery
}

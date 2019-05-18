package basic

import (
	"fmt"
	brainException "gitlab.com/iotTracker/brain/exception"
	"gitlab.com/iotTracker/brain/tracker/zx303"
	zx303StatusReadingReportGenerator "gitlab.com/iotTracker/brain/tracker/zx303/reading/status/report/generator"
)

type generator struct {
}

func New() zx303StatusReadingReportGenerator.Generator {
	return &generator{}
}

func (g *generator) ValidateBatteryRequest(request *zx303StatusReadingReportGenerator.BatteryRequest) error {
	reasonsInvalid := make([]string, 0)

	if request.Claims == nil {
		reasonsInvalid = append(reasonsInvalid, "claims are nil")
	}

	if request.ZX303TrackerIdentifier == nil {
		reasonsInvalid = append(reasonsInvalid, "claims are nil")
	} else if !zx303.IsValidIdentifier(request.ZX303TrackerIdentifier) {
		reasonsInvalid = append(reasonsInvalid, "invalid identifier")
	}

	if len(reasonsInvalid) > 0 {
		return brainException.RequestInvalid{Reasons: reasonsInvalid}
	}
	return nil
}

func (g *generator) Battery(request *zx303StatusReadingReportGenerator.BatteryRequest) (*zx303StatusReadingReportGenerator.BatteryResponse, error) {
	if err := g.ValidateBatteryRequest(request); err != nil {
		return nil, err
	}

	fmt.Println("running batter report!")

	return &zx303StatusReadingReportGenerator.BatteryResponse{}, nil
}

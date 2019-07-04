package basic

import (
	brainException "github.com/iot-my-world/brain/internal/exception"
	"github.com/iot-my-world/brain/internal/log"
	"github.com/iot-my-world/brain/pkg/search/criterion"
	exactText "github.com/iot-my-world/brain/pkg/search/criterion/exact/text"
	dateRange "github.com/iot-my-world/brain/pkg/search/criterion/range/date"
	zx3032 "github.com/iot-my-world/brain/pkg/tracker/zx303"
	recordHandler2 "github.com/iot-my-world/brain/pkg/tracker/zx303/reading/status/recordHandler"
	"github.com/iot-my-world/brain/pkg/tracker/zx303/reading/status/report"
	generator2 "github.com/iot-my-world/brain/pkg/tracker/zx303/reading/status/report/generator"
	"github.com/iot-my-world/brain/pkg/tracker/zx303/reading/status/report/generator/exception"
	"github.com/iot-my-world/brain/pkg/tracker/zx303/recordHandler"
)

type generator struct {
	zx303StatusReadingRecordHandler *recordHandler2.RecordHandler
	zx303RecordHandler              *recordHandler.RecordHandler
}

func New(
	zx303StatusReadingRecordHandler *recordHandler2.RecordHandler,
	zx303RecordHandler *recordHandler.RecordHandler,
) generator2.Generator {
	return &generator{
		zx303StatusReadingRecordHandler: zx303StatusReadingRecordHandler,
		zx303RecordHandler:              zx303RecordHandler,
	}
}

func (g *generator) ValidateBatteryReportRequest(request *generator2.BatteryReportRequest) error {
	reasonsInvalid := make([]string, 0)

	if request.Claims == nil {
		reasonsInvalid = append(reasonsInvalid, "claims are nil")
	}

	if request.ZX303TrackerIdentifier == nil {
		reasonsInvalid = append(reasonsInvalid, "ZX303TrackerIdentifier is nil")
	} else if !zx3032.IsValidIdentifier(request.ZX303TrackerIdentifier) {
		reasonsInvalid = append(reasonsInvalid, "invalid ZX303TrackerIdentifier")
	}

	if len(reasonsInvalid) > 0 {
		return brainException.RequestInvalid{Reasons: reasonsInvalid}
	}
	return nil
}

func (g *generator) BatteryReport(request *generator2.BatteryReportRequest) (*generator2.BatteryReportResponse, error) {
	if err := g.ValidateBatteryReportRequest(request); err != nil {
		return nil, err
	}

	// retrieve the device
	retrieveResponse, err := g.zx303RecordHandler.Retrieve(&recordHandler.RetrieveRequest{
		Claims:     request.Claims,
		Identifier: request.ZX303TrackerIdentifier,
	})
	if err != nil {
		err = exception.BatteryReportGeneration{
			Reasons: []string{
				"retrieving device",
				err.Error(),
			},
		}
		log.Error(err.Error())
		return nil, err
	}

	// collect all of the status readings for this device
	readingCollectResponse, err := g.zx303StatusReadingRecordHandler.Collect(&recordHandler2.CollectRequest{
		Claims: request.Claims,
		Criteria: []criterion.Criterion{
			exactText.Criterion{
				Field: "deviceId.id",
				Text:  retrieveResponse.ZX303.Id,
			},
			dateRange.Criterion{
				Field: "timeStamp",
				StartDate: dateRange.RangeValue{
					Date:      request.StartDate,
					Inclusive: true,
					Ignore:    false,
				},
				EndDate: dateRange.RangeValue{
					Date:      request.EndDate,
					Inclusive: true,
					Ignore:    false,
				},
			},
		},
		//Query:    query.Query{},
	})
	if err != nil {
		err = exception.BatteryReportGeneration{
			Reasons: []string{
				"collecting readings",
				err.Error(),
			},
		}
		log.Error(err.Error())
		return nil, err
	}

	batteryReport := report.Battery{
		Readings: make([]report.BatteryReading, 0),
	}

	for _, reading := range readingCollectResponse.Records {
		batteryReport.Readings = append(batteryReport.Readings, report.BatteryReading{
			Percentage: reading.BatteryPercentage,
			Timestamp:  reading.Timestamp,
		})
	}

	return &generator2.BatteryReportResponse{
		Report: batteryReport,
	}, nil
}

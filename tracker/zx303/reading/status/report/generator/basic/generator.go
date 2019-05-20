package basic

import (
	brainException "gitlab.com/iotTracker/brain/exception"
	"gitlab.com/iotTracker/brain/log"
	"gitlab.com/iotTracker/brain/search/criterion"
	exactText "gitlab.com/iotTracker/brain/search/criterion/exact/text"
	dateRange "gitlab.com/iotTracker/brain/search/criterion/range/date"
	"gitlab.com/iotTracker/brain/tracker/zx303"
	zx303StatusReadingRecordHandler "gitlab.com/iotTracker/brain/tracker/zx303/reading/status/recordHandler"
	zx303StatusReadingReport "gitlab.com/iotTracker/brain/tracker/zx303/reading/status/report"
	zx303StatusReadingReportGenerator "gitlab.com/iotTracker/brain/tracker/zx303/reading/status/report/generator"
	zx303StatusReadingReportGeneratorException "gitlab.com/iotTracker/brain/tracker/zx303/reading/status/report/generator/exception"
	zx303RecordHandler "gitlab.com/iotTracker/brain/tracker/zx303/recordHandler"
)

type generator struct {
	zx303StatusReadingRecordHandler *zx303StatusReadingRecordHandler.RecordHandler
	zx303RecordHandler              *zx303RecordHandler.RecordHandler
}

func New(
	zx303StatusReadingRecordHandler *zx303StatusReadingRecordHandler.RecordHandler,
	zx303RecordHandler *zx303RecordHandler.RecordHandler,
) zx303StatusReadingReportGenerator.Generator {
	return &generator{
		zx303StatusReadingRecordHandler: zx303StatusReadingRecordHandler,
		zx303RecordHandler:              zx303RecordHandler,
	}
}

func (g *generator) ValidateBatteryReportRequest(request *zx303StatusReadingReportGenerator.BatteryReportRequest) error {
	reasonsInvalid := make([]string, 0)

	if request.Claims == nil {
		reasonsInvalid = append(reasonsInvalid, "claims are nil")
	}

	if request.ZX303TrackerIdentifier == nil {
		reasonsInvalid = append(reasonsInvalid, "ZX303TrackerIdentifier is nil")
	} else if !zx303.IsValidIdentifier(request.ZX303TrackerIdentifier) {
		reasonsInvalid = append(reasonsInvalid, "invalid ZX303TrackerIdentifier")
	}

	if len(reasonsInvalid) > 0 {
		return brainException.RequestInvalid{Reasons: reasonsInvalid}
	}
	return nil
}

func (g *generator) BatteryReport(request *zx303StatusReadingReportGenerator.BatteryReportRequest) (*zx303StatusReadingReportGenerator.BatteryReportResponse, error) {
	if err := g.ValidateBatteryReportRequest(request); err != nil {
		return nil, err
	}

	// retrieve the device
	retrieveResponse, err := g.zx303RecordHandler.Retrieve(&zx303RecordHandler.RetrieveRequest{
		Claims:     request.Claims,
		Identifier: request.ZX303TrackerIdentifier,
	})
	if err != nil {
		err = zx303StatusReadingReportGeneratorException.BatteryReportGeneration{
			Reasons: []string{
				"retrieving device",
				err.Error(),
			},
		}
		log.Error(err.Error())
		return nil, err
	}

	// collect all of the status readings for this device
	readingCollectResponse, err := g.zx303StatusReadingRecordHandler.Collect(&zx303StatusReadingRecordHandler.CollectRequest{
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
		err = zx303StatusReadingReportGeneratorException.BatteryReportGeneration{
			Reasons: []string{
				"collecting readings",
				err.Error(),
			},
		}
		log.Error(err.Error())
		return nil, err
	}

	batteryReport := zx303StatusReadingReport.Battery{
		Readings: make([]zx303StatusReadingReport.BatteryReading, 0),
	}

	for _, reading := range readingCollectResponse.Records {
		batteryReport.Readings = append(batteryReport.Readings, zx303StatusReadingReport.BatteryReading{
			Percentage: reading.BatteryPercentage,
			Timestamp:  reading.Timestamp,
		})
	}

	return &zx303StatusReadingReportGenerator.BatteryReportResponse{
		Report: batteryReport,
	}, nil
}

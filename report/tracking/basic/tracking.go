package basic

import (
	brainException "gitlab.com/iotTracker/brain/exception"
	companyRecordHandler "gitlab.com/iotTracker/brain/party/company/recordHandler"
	clientRecordHandler "gitlab.com/iotTracker/brain/party/client/recordHandler"
	readingRecordHandler "gitlab.com/iotTracker/brain/tracker/reading/recordHandler"
	trackingReport "gitlab.com/iotTracker/brain/report/tracking"
)

type basicTrackingReport struct {
	companyRecordHandler companyRecordHandler.RecordHandler
	clientRecordHandler  clientRecordHandler.RecordHandler
	readingRecordHandler readingRecordHandler.RecordHandler
}

func New(
	companyRecordHandler companyRecordHandler.RecordHandler,
	clientRecordHandler clientRecordHandler.RecordHandler,
	readingRecordHandler readingRecordHandler.RecordHandler,
) *basicTrackingReport {
	return &basicTrackingReport{
		companyRecordHandler: companyRecordHandler,
		clientRecordHandler:  clientRecordHandler,
		readingRecordHandler: readingRecordHandler,
	}
}

func (btr *basicTrackingReport) ValidateLiveRequest(request *trackingReport.LiveRequest) error {
	reasonsInvalid := make([]string, 0)

	if len(reasonsInvalid) > 0 {
		return brainException.RequestInvalid{Reasons: reasonsInvalid}
	}

	return nil
}

func (btr *basicTrackingReport) Live(request *trackingReport.LiveRequest, response *trackingReport.LiveResponse) error {
	if err := btr.ValidateLiveRequest(request); err != nil {
		return err
	}

	return nil
}

func (btr *basicTrackingReport) ValidateHistoricalRequest(request *trackingReport.HistoricalRequest) error {
	reasonsInvalid := make([]string, 0)

	if len(reasonsInvalid) > 0 {
		return brainException.RequestInvalid{Reasons: reasonsInvalid}
	}

	return nil
}

func (btr *basicTrackingReport) Historical(request *trackingReport.HistoricalRequest, response *trackingReport.HistoricalResponse) error {
	if err := btr.ValidateHistoricalRequest(request); err != nil {
		return err
	}

	return nil
}

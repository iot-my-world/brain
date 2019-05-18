package basic

import (
	brainException "gitlab.com/iotTracker/brain/exception"
	partyAdministrator "gitlab.com/iotTracker/brain/party/administrator"
	trackingReport "gitlab.com/iotTracker/brain/report/tracking"
	trackingReportException "gitlab.com/iotTracker/brain/report/tracking/exception"
	"gitlab.com/iotTracker/brain/search/criterion"
	exactTextCriterion "gitlab.com/iotTracker/brain/search/criterion/exact/text"
	textListCriterion "gitlab.com/iotTracker/brain/search/criterion/list/text"
	"gitlab.com/iotTracker/brain/search/criterion/or"
	"gitlab.com/iotTracker/brain/search/query"
	zx303GPSReading "gitlab.com/iotTracker/brain/tracker/zx303/reading/gps"
	zx303TrackerReadingRecordHandler "gitlab.com/iotTracker/brain/tracker/zx303/reading/gps/recordHandler"
	zx303TrackerRecordHandler "gitlab.com/iotTracker/brain/tracker/zx303/recordHandler"
)

type basicTrackingReport struct {
	partyAdministrator               partyAdministrator.Administrator
	zx303TrackerReadingRecordHandler *zx303TrackerReadingRecordHandler.RecordHandler
	zx303TrackerRecordHandler        *zx303TrackerRecordHandler.RecordHandler
}

func New(
	partyAdministrator partyAdministrator.Administrator,
	zx303TrackerReadingRecordHandler *zx303TrackerReadingRecordHandler.RecordHandler,
	zx303TrackerRecordHandler *zx303TrackerRecordHandler.RecordHandler,
) trackingReport.Report {
	return &basicTrackingReport{
		partyAdministrator:               partyAdministrator,
		zx303TrackerReadingRecordHandler: zx303TrackerReadingRecordHandler,
		zx303TrackerRecordHandler:        zx303TrackerRecordHandler,
	}
}

func (btr *basicTrackingReport) ValidateLiveRequest(request *trackingReport.LiveRequest) error {
	reasonsInvalid := make([]string, 0)

	// confirm that the claims are not nil
	if request.Claims == nil {
		reasonsInvalid = append(reasonsInvalid, "claims are nil")
	}

	// confirm that all of the identifiers are valid
	for idIdx := range request.PartyIdentifiers {
		if err := request.PartyIdentifiers[idIdx].IsValid(); err != nil {
			reasonsInvalid = append(reasonsInvalid, "invalid party identifier"+err.Error())
			break
		}
	}

	if len(reasonsInvalid) > 0 {
		return brainException.RequestInvalid{Reasons: reasonsInvalid}
	}

	return nil
}

func (btr *basicTrackingReport) Live(request *trackingReport.LiveRequest) (*trackingReport.LiveResponse, error) {
	if err := btr.ValidateLiveRequest(request); err != nil {
		return nil, err
	}

	// records to return
	zx303GPSLiveReportReadings := make([]zx303GPSReading.Reading, 0)

	// query for collecting only latest reading
	collectQuery := query.Query{
		Limit:  1,
		Offset: 0,
		Order:  []query.SortOrder{query.SortOrderDescending},
		SortBy: []string{"timeStamp"},
	}

	// retrieve each party with provided identifiers
	for _, partyIdentifier := range request.PartyIdentifiers {
		retrieveResponse, err := btr.partyAdministrator.RetrieveParty(&partyAdministrator.RetrievePartyRequest{
			Claims:     request.Claims,
			Identifier: partyIdentifier.PartyIdIdentifier,
			PartyType:  partyIdentifier.PartyType,
		})
		if err != nil {
			return nil, trackingReportException.RetrievingParty{Reasons: []string{string(partyIdentifier.PartyType), err.Error()}}
		}

		// criterion to collect all devices either owned by or assigned to the party
		collectCriterion := or.Criterion{
			Criteria: []criterion.Criterion{
				textListCriterion.Criterion{
					Field: "ownerId.id",
					List:  []string{retrieveResponse.Party.Details().PartyId.Id},
				},
				textListCriterion.Criterion{
					Field: "assignedId.id",
					List:  []string{retrieveResponse.Party.Details().PartyId.Id},
				},
			},
		}

		// collect all the zx303 devices
		zx303TrackerCollectResponse, err := btr.zx303TrackerRecordHandler.Collect(&zx303TrackerRecordHandler.CollectRequest{
			Claims:   request.Claims,
			Criteria: []criterion.Criterion{collectCriterion},
			// Query: left blank to collect all. i.e. no limit
		})
		if err != nil {
			return nil, trackingReportException.CollectingDevices{Reasons: []string{"zx303 devices", err.Error()}}
		}

		// collect the last reading associated with each of these devices
		for devIdx := range zx303TrackerCollectResponse.Records {
			// exact text criterion for this device
			deviceIDExactTextCriterion := exactTextCriterion.Criterion{
				Field: "deviceId.id",
				Text:  zx303TrackerCollectResponse.Records[devIdx].Id,
			}

			// collect the latest reading for this device
			readingCollectResponse, err := btr.zx303TrackerReadingRecordHandler.Collect(&zx303TrackerReadingRecordHandler.CollectRequest{
				Claims:   request.Claims,
				Query:    collectQuery,
				Criteria: []criterion.Criterion{deviceIDExactTextCriterion},
			})
			if err != nil {
				return nil, trackingReportException.CollectingReadings{Reasons: []string{"tk102 device readings", err.Error()}}
			}
			// if any readings have been collected for this device
			if len(readingCollectResponse.Records) > 0 {
				if len(zx303GPSLiveReportReadings) == 0 {
					// if noting has been added to the live report readings yet, add it now
					zx303GPSLiveReportReadings = append(zx303GPSLiveReportReadings, readingCollectResponse.Records[0])
				} else {
					// otherwise check if the reading has not yet been added
					for readingIdx := range zx303GPSLiveReportReadings {
						if zx303GPSLiveReportReadings[readingIdx].Id == readingCollectResponse.Records[0].Id {
							// it has already been added, break
							break
						}
						if readingIdx == len(zx303GPSLiveReportReadings)-1 {
							// it has not been added, add it now
							zx303GPSLiveReportReadings = append(zx303GPSLiveReportReadings, readingCollectResponse.Records[0])
						}
					}
				}
			}
		}
	}

	return &trackingReport.LiveResponse{
		ZX303TrackerGPSReadings: zx303GPSLiveReportReadings,
	}, nil
}

func (btr *basicTrackingReport) ValidateHistoricalRequest(request *trackingReport.HistoricalRequest) error {
	reasonsInvalid := make([]string, 0)

	if len(reasonsInvalid) > 0 {
		return brainException.RequestInvalid{Reasons: reasonsInvalid}
	}

	return nil
}

func (btr *basicTrackingReport) Historical(request *trackingReport.HistoricalRequest) (*trackingReport.HistoricalResponse, error) {
	if err := btr.ValidateHistoricalRequest(request); err != nil {
		return nil, err
	}

	return &trackingReport.HistoricalResponse{
		//Readings: make([]reading.Reading, 0),
	}, nil
}

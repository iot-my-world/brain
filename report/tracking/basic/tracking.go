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
	"gitlab.com/iotTracker/brain/tracker/reading"
	readingRecordHandler "gitlab.com/iotTracker/brain/tracker/reading/recordHandler"
	tk102DeviceRecordHandler "gitlab.com/iotTracker/brain/tracker/tk102/recordHandler"
)

type basicTrackingReport struct {
	partyAdministrator       partyAdministrator.Administrator
	readingRecordHandler     readingRecordHandler.RecordHandler
	tk102DeviceRecordHandler tk102DeviceRecordHandler.RecordHandler
}

// New basic tracking report
func New(
	partyAdministrator partyAdministrator.Administrator,
	readingRecordHandler readingRecordHandler.RecordHandler,
	tk102DeviceRecordHandler tk102DeviceRecordHandler.RecordHandler,
) trackingReport.Report {
	return &basicTrackingReport{
		partyAdministrator:       partyAdministrator,
		readingRecordHandler:     readingRecordHandler,
		tk102DeviceRecordHandler: tk102DeviceRecordHandler,
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
	liveReportReadings := make([]reading.Reading, 0)

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

		// collect all the tk102 devices
		tk102DeviceCollectResponse, err := btr.tk102DeviceRecordHandler.Collect(&tk102DeviceRecordHandler.CollectRequest{
			Claims:   request.Claims,
			Criteria: []criterion.Criterion{collectCriterion},
			// Query: left blank to collect all. i.e. no limit
		})
		if err != nil {
			return nil, trackingReportException.CollectingDevices{Reasons: []string{"tk102 devices", err.Error()}}
		}

		// collect the last reading associated with each of these devices
		for devIdx := range tk102DeviceCollectResponse.Records {
			// exact text criterion for this device
			deviceIDExactTextCriterion := exactTextCriterion.Criterion{
				Field: "deviceId.id",
				Text:  tk102DeviceCollectResponse.Records[devIdx].Id,
			}

			// collect the latest reading for this device
			readingCollectResponse, err := btr.readingRecordHandler.Collect(&readingRecordHandler.CollectRequest{
				Claims:   request.Claims,
				Query:    collectQuery,
				Criteria: []criterion.Criterion{deviceIDExactTextCriterion},
			})
			if err != nil {
				return nil, trackingReportException.CollectingReadings{Reasons: []string{"tk102 device readings", err.Error()}}
			}
			// if any readings have been collected for this device
			if len(readingCollectResponse.Records) > 0 {
				if len(liveReportReadings) == 0 {
					// if noting has been added to the live report readings yet, add it now
					liveReportReadings = append(liveReportReadings, readingCollectResponse.Records[0])
				} else {
					// otherwise check if the reading has not yet been added
					for readingIdx := range liveReportReadings {
						if liveReportReadings[readingIdx].Id == readingCollectResponse.Records[0].Id {
							// it has already been added, break
							break
						}
						if readingIdx == len(liveReportReadings)-1 {
							// it has not been added, add it now
							liveReportReadings = append(liveReportReadings, readingCollectResponse.Records[0])
						}
					}
				}
			}
		}
	}

	return &trackingReport.LiveResponse{
		Readings: liveReportReadings,
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
		Readings: make([]reading.Reading, 0),
	}, nil
}

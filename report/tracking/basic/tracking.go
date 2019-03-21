package basic

import (
	brainException "gitlab.com/iotTracker/brain/exception"
	"gitlab.com/iotTracker/brain/party/client"
	clientRecordHandler "gitlab.com/iotTracker/brain/party/client/recordHandler"
	"gitlab.com/iotTracker/brain/party/company"
	companyRecordHandler "gitlab.com/iotTracker/brain/party/company/recordHandler"
	"gitlab.com/iotTracker/brain/party/system"
	systemRecordHandler "gitlab.com/iotTracker/brain/party/system/recordHandler"
	trackingReport "gitlab.com/iotTracker/brain/report/tracking"
	trackingReportException "gitlab.com/iotTracker/brain/report/tracking/exception"
	"gitlab.com/iotTracker/brain/search/criterion"
	exactTextCriterion "gitlab.com/iotTracker/brain/search/criterion/exact/text"
	textListCriterion "gitlab.com/iotTracker/brain/search/criterion/list/text"
	"gitlab.com/iotTracker/brain/search/criterion/or"
	"gitlab.com/iotTracker/brain/search/query"
	tk102DeviceRecordHandler "gitlab.com/iotTracker/brain/tracker/device/tk102/recordHandler"
	"gitlab.com/iotTracker/brain/tracker/reading"
	readingRecordHandler "gitlab.com/iotTracker/brain/tracker/reading/recordHandler"
)

type basicTrackingReport struct {
	systemRecordHandler      systemRecordHandler.RecordHandler
	companyRecordHandler     companyRecordHandler.RecordHandler
	clientRecordHandler      clientRecordHandler.RecordHandler
	readingRecordHandler     readingRecordHandler.RecordHandler
	tk102DeviceRecordHandler tk102DeviceRecordHandler.RecordHandler
}

// New basic tracking report
func New(
	systemRecordHandler systemRecordHandler.RecordHandler,
	companyRecordHandler companyRecordHandler.RecordHandler,
	clientRecordHandler clientRecordHandler.RecordHandler,
	readingRecordHandler readingRecordHandler.RecordHandler,
	tk102DeviceRecordHandler tk102DeviceRecordHandler.RecordHandler,
) *basicTrackingReport {
	return &basicTrackingReport{
		systemRecordHandler:      systemRecordHandler,
		companyRecordHandler:     companyRecordHandler,
		clientRecordHandler:      clientRecordHandler,
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

	// confirm that none of the identifiers are nil and that they are valid
	for idIdx := range request.SystemIdentifiers {
		if request.SystemIdentifiers[idIdx] == nil {
			reasonsInvalid = append(reasonsInvalid, "nil system identifier")
			break
		} else if !system.IsValidIdentifier(request.SystemIdentifiers[idIdx]) {
			reasonsInvalid = append(reasonsInvalid, "invalid system identifier")
			break
		}
	}
	// confirm that none of the identifiers are nil and that they are valid
	for idIdx := range request.ClientIdentifiers {
		if request.ClientIdentifiers[idIdx] == nil {
			reasonsInvalid = append(reasonsInvalid, "nil client identifier")
			break
		} else if !client.IsValidIdentifier(request.ClientIdentifiers[idIdx]) {
			reasonsInvalid = append(reasonsInvalid, "invalid client identifier")
			break
		}
	}
	for idIdx := range request.CompanyIdentifiers {
		if request.CompanyIdentifiers[idIdx] == nil {
			reasonsInvalid = append(reasonsInvalid, "nil company identifier")
		} else if !company.IsValidIdentifier(request.ClientIdentifiers[idIdx]) {
			reasonsInvalid = append(reasonsInvalid, "invalid company identifier")
			break
		}
	}

	if len(reasonsInvalid) > 0 {
		return brainException.RequestInvalid{Reasons: reasonsInvalid}
	}

	return nil
}

func (btr *basicTrackingReport) Live(request *trackingReport.LiveRequest, response *trackingReport.LiveResponse) error {
	if err := btr.ValidateLiveRequest(request); err != nil {
		return err
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

	// For each system identifier provided
	for _, id := range request.SystemIdentifiers {
		// retrieve the system
		systemRetrieveResponse := systemRecordHandler.RetrieveResponse{}
		if err := btr.systemRecordHandler.Retrieve(&systemRecordHandler.RetrieveRequest{
			Identifier: id,
			Claims:     request.Claims,
		}, &systemRetrieveResponse); err != nil {
			return trackingReportException.RetrievingSystem{Reasons: []string{err.Error()}}
		}

		// criterion to collect all devices either owned by or assigned to the system
		collectCriterion := or.Criterion{
			Criteria: []criterion.Criterion{
				textListCriterion.Criterion{
					Field: "ownerId.id",
					List:  []string{systemRetrieveResponse.System.Id},
				},
				textListCriterion.Criterion{
					Field: "assignedId.id",
					List:  []string{systemRetrieveResponse.System.Id},
				},
			},
		}

		// collect all the tk102 devices
		tk102DeviceCollectResponse := tk102DeviceRecordHandler.CollectResponse{}
		if err := btr.tk102DeviceRecordHandler.Collect(&tk102DeviceRecordHandler.CollectRequest{
			Claims:   request.Claims,
			Criteria: []criterion.Criterion{collectCriterion},
			// Query: left blank to collect all. i.e. no limit
		}, &tk102DeviceCollectResponse); err != nil {
			return trackingReportException.CollectingDevices{Reasons: []string{"tk102 devices", err.Error()}}
		}

		// collect the last reading associated with each of these devices
		for devIdx := range tk102DeviceCollectResponse.Records {
			// exact text criterion for this device
			deviceIDExactTextCriterion := exactTextCriterion.Criterion{
				Field: "deviceId.id",
				Text:  tk102DeviceCollectResponse.Records[devIdx].Id,
			}

			// collect the latest reading for this device
			readingCollectResponse := readingRecordHandler.CollectResponse{}
			if err := btr.readingRecordHandler.Collect(&readingRecordHandler.CollectRequest{
				Claims:   request.Claims,
				Query:    collectQuery,
				Criteria: []criterion.Criterion{deviceIDExactTextCriterion},
			}, &readingCollectResponse); err != nil {
				return trackingReportException.CollectingReadings{Reasons: []string{"tk102 device readings", err.Error()}}
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

	// For each company identifier provided
	for _, id := range request.CompanyIdentifiers {
		// retrieve the company
		companyRetrieveResponse := companyRecordHandler.RetrieveResponse{}
		if err := btr.companyRecordHandler.Retrieve(&companyRecordHandler.RetrieveRequest{
			Identifier: id,
			Claims:     request.Claims,
		}, &companyRetrieveResponse); err != nil {
			return trackingReportException.RetrievingCompany{Reasons: []string{err.Error()}}
		}

		// criterion to collect all devices either owned by or assigned to the company
		collectCriterion := or.Criterion{
			Criteria: []criterion.Criterion{
				textListCriterion.Criterion{
					Field: "ownerId.id",
					List:  []string{companyRetrieveResponse.Company.Id},
				},
				textListCriterion.Criterion{
					Field: "assignedId.id",
					List:  []string{companyRetrieveResponse.Company.Id},
				},
			},
		}

		// collect all the tk102 devices
		tk102DeviceCollectResponse := tk102DeviceRecordHandler.CollectResponse{}
		if err := btr.tk102DeviceRecordHandler.Collect(&tk102DeviceRecordHandler.CollectRequest{
			Claims:   request.Claims,
			Criteria: []criterion.Criterion{collectCriterion},
			// Query: left blank to collect all. i.e. no limit
		}, &tk102DeviceCollectResponse); err != nil {
			return trackingReportException.CollectingDevices{Reasons: []string{"tk102 devices", err.Error()}}
		}

		// collect the last reading associated with each of these devices
		for devIdx := range tk102DeviceCollectResponse.Records {
			// exact text criterion for this device
			deviceIDExactTextCriterion := exactTextCriterion.Criterion{
				Field: "deviceId.id",
				Text:  tk102DeviceCollectResponse.Records[devIdx].Id,
			}

			// collect the latest reading for this device
			readingCollectResponse := readingRecordHandler.CollectResponse{}
			if err := btr.readingRecordHandler.Collect(&readingRecordHandler.CollectRequest{
				Claims:   request.Claims,
				Query:    collectQuery,
				Criteria: []criterion.Criterion{deviceIDExactTextCriterion},
			}, &readingCollectResponse); err != nil {
				return trackingReportException.CollectingReadings{Reasons: []string{"tk102 device readings", err.Error()}}
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

	// For each client identifier provided
	for _, id := range request.ClientIdentifiers {
		// retrieve the client
		clientRetrieveResponse, err := btr.clientRecordHandler.Retrieve(&clientRecordHandler.RetrieveRequest{
			Identifier: id,
			Claims:     request.Claims,
		})
		if err != nil {
			return trackingReportException.RetrievingClient{Reasons: []string{err.Error()}}
		}

		// criterion to collect all devices either owned by or assigned to the client
		collectCriterion := or.Criterion{
			Criteria: []criterion.Criterion{
				textListCriterion.Criterion{
					Field: "ownerId.id",
					List:  []string{clientRetrieveResponse.Client.Id},
				},
				textListCriterion.Criterion{
					Field: "assignedId.id",
					List:  []string{clientRetrieveResponse.Client.Id},
				},
			},
		}

		// collect all the tk102 devices
		tk102DeviceCollectResponse := tk102DeviceRecordHandler.CollectResponse{}
		if err := btr.tk102DeviceRecordHandler.Collect(&tk102DeviceRecordHandler.CollectRequest{
			Claims:   request.Claims,
			Criteria: []criterion.Criterion{collectCriterion},
			// Query: left blank to collect all. i.e. no limit
		}, &tk102DeviceCollectResponse); err != nil {
			return trackingReportException.CollectingDevices{Reasons: []string{"tk102 devices", err.Error()}}
		}

		// collect the last reading associated with each of these devices
		for devIdx := range tk102DeviceCollectResponse.Records {
			// exact text criterion for this device
			deviceIDExactTextCriterion := exactTextCriterion.Criterion{
				Field: "deviceId.id",
				Text:  tk102DeviceCollectResponse.Records[devIdx].Id,
			}

			// collect the latest reading for this device
			readingCollectResponse := readingRecordHandler.CollectResponse{}
			if err := btr.readingRecordHandler.Collect(&readingRecordHandler.CollectRequest{
				Claims:   request.Claims,
				Query:    collectQuery,
				Criteria: []criterion.Criterion{deviceIDExactTextCriterion},
			}, &readingCollectResponse); err != nil {
				return trackingReportException.CollectingReadings{Reasons: []string{"tk102 device readings", err.Error()}}
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

	response.Readings = liveReportReadings

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

	response.Readings = make([]reading.Reading, 0)

	return nil
}

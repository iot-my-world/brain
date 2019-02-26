package system

import (
	"gitlab.com/iotTracker/brain/party/system"
	systemRecordHandler "gitlab.com/iotTracker/brain/party/system/recordHandler"
	"gitlab.com/iotTracker/brain/search/wrappedIdentifier"
	"net/http"
	"gitlab.com/iotTracker/brain/search/wrappedCriterion"
	"gitlab.com/iotTracker/brain/search/criterion"
	"gitlab.com/iotTracker/brain/search/query"
)

type adaptor struct {
	RecordHandler systemRecordHandler.RecordHandler
}

func New(recordHandler systemRecordHandler.RecordHandler) *adaptor {
	return &adaptor{
		RecordHandler: recordHandler,
	}
}

type RetrieveRequest struct {
	Identifier wrappedIdentifier.WrappedIdentifier `json:"identifier"`
}

type RetrieveResponse struct {
	System system.System `json:"system" bson:"system"`
}

func (s *adaptor) Retrieve(r *http.Request, request *RetrieveRequest, response *RetrieveResponse) error {
	id, err := request.Identifier.UnWrap()
	if err != nil {
		return err
	}

	retrieveSystemResponse := systemRecordHandler.RetrieveResponse{}
	if err := s.RecordHandler.Retrieve(
		&systemRecordHandler.RetrieveRequest{
			Identifier: id,
		},
		&retrieveSystemResponse); err != nil {
		return err
	}

	response.System = retrieveSystemResponse.System

	return nil
}

type CollectRequest struct {
	Criteria []wrappedCriterion.WrappedCriterion `json:"criteria"`
	Query    query.Query                         `json:"query"`
}

type CollectResponse struct {
	Records []system.System `json:"records"`
	Total   int             `json:"total"`
}

func (s *adaptor) Collect(r *http.Request, request *CollectRequest, response *CollectResponse) error {
	// unwrap criteria
	criteria := make([]criterion.Criterion, 0)
	for criterionIdx := range request.Criteria {
		if c, err := request.Criteria[criterionIdx].UnWrap(); err == nil {
			criteria = append(criteria, c)
		} else {
			return err
		}
	}

	collectSystemResponse := systemRecordHandler.CollectResponse{}
	if err := s.RecordHandler.Collect(&systemRecordHandler.CollectRequest{
		Criteria: criteria,
		Query:    request.Query,
	},
		&collectSystemResponse); err != nil {
		return err
	}

	response.Records = collectSystemResponse.Records
	response.Total = collectSystemResponse.Total
	return nil
}

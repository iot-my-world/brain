package jsonRpc

import (
	readingRecordHandler "gitlab.com/iotTracker/brain/tracker/reading/recordHandler"
	"gitlab.com/iotTracker/brain/search/wrappedCriterion"
	"gitlab.com/iotTracker/brain/tracker/reading"
	"net/http"
	"gitlab.com/iotTracker/brain/search/criterion"
	"gitlab.com/iotTracker/brain/search/query"
)

type adaptor struct {
	RecordHandler readingRecordHandler.RecordHandler
}

func New(recordHandler readingRecordHandler.RecordHandler) *adaptor {
	return &adaptor{
		RecordHandler: recordHandler,
	}
}

type CollectRequest struct {
	Criteria []wrappedCriterion.WrappedCriterion `json:"criteria"`
	Query    query.Query                         `json:"query"`
}

type CollectResponse struct {
	Records []reading.Reading `json:"records"`
	Total   int               `json:"total"`
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

	collectReadingResponse := readingRecordHandler.CollectResponse{}
	if err := s.RecordHandler.Collect(&readingRecordHandler.CollectRequest{
		Criteria: criteria,
		Query:    request.Query,
	},
		&collectReadingResponse); err != nil {
		return err
	}

	response.Records = collectReadingResponse.Records
	response.Total = collectReadingResponse.Total
	return nil
}

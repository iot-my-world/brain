package jsonRpc

import (
	"gitlab.com/iotTracker/brain/log"
	"gitlab.com/iotTracker/brain/search/criterion"
	wrappedCriterion "gitlab.com/iotTracker/brain/search/criterion/wrapped"
	"gitlab.com/iotTracker/brain/search/query"
	wrappedClaims "gitlab.com/iotTracker/brain/security/claims/wrapped"
	"gitlab.com/iotTracker/brain/tracker/tk102/reading"
	readingRecordHandler "gitlab.com/iotTracker/brain/tracker/tk102/reading/recordHandler"
	"net/http"
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
	Criteria []wrappedCriterion.Wrapped `json:"criteria"`
	Query    query.Query                `json:"query"`
}

type CollectResponse struct {
	Records []reading.Reading `json:"records"`
	Total   int               `json:"total"`
}

func (s *adaptor) Collect(r *http.Request, request *CollectRequest, response *CollectResponse) error {
	claims, err := wrappedClaims.UnwrapClaimsFromContext(r)
	if err != nil {
		log.Warn(err.Error())
		return err
	}

	criteria := make([]criterion.Criterion, 0)
	for criterionIdx := range request.Criteria {
		if c, err := request.Criteria[criterionIdx].UnWrap(); err == nil {
			criteria = append(criteria, c)
		} else {
			return err
		}
	}

	collectReadingResponse, err := s.RecordHandler.Collect(&readingRecordHandler.CollectRequest{
		Claims:   claims,
		Criteria: criteria,
		Query:    request.Query,
	})
	if err != nil {
		return err
	}

	response.Records = collectReadingResponse.Records
	response.Total = collectReadingResponse.Total
	return nil
}
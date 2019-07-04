package jsonRpc

import (
	"github.com/iot-my-world/brain/log"
	"github.com/iot-my-world/brain/pkg/search/criterion"
	wrappedCriterion "github.com/iot-my-world/brain/pkg/search/criterion/wrapped"
	"github.com/iot-my-world/brain/pkg/search/query"
	reading2 "github.com/iot-my-world/brain/pkg/tracker/tk102/reading"
	"github.com/iot-my-world/brain/pkg/tracker/tk102/reading/recordHandler"
	wrappedClaims "github.com/iot-my-world/brain/security/claims/wrapped"
	"net/http"
)

type adaptor struct {
	RecordHandler recordHandler.RecordHandler
}

func New(recordHandler recordHandler.RecordHandler) *adaptor {
	return &adaptor{
		RecordHandler: recordHandler,
	}
}

type CollectRequest struct {
	Criteria []wrappedCriterion.Wrapped `json:"criteria"`
	Query    query.Query                `json:"query"`
}

type CollectResponse struct {
	Records []reading2.Reading `json:"records"`
	Total   int                `json:"total"`
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

	collectReadingResponse, err := s.RecordHandler.Collect(&recordHandler.CollectRequest{
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

package jsonRpc

import (
	"github.com/iot-my-world/brain/internal/log"
	"github.com/iot-my-world/brain/pkg/search/criterion"
	wrappedCriterion "github.com/iot-my-world/brain/pkg/search/criterion/wrapped"
	wrappedIdentifier "github.com/iot-my-world/brain/pkg/search/identifier/wrapped"
	"github.com/iot-my-world/brain/pkg/search/query"
	wrappedClaims "github.com/iot-my-world/brain/pkg/security/claims/wrapped"
	"github.com/iot-my-world/brain/pkg/tracker/zx303/reading/status"
	"github.com/iot-my-world/brain/pkg/tracker/zx303/reading/status/recordHandler"
	"net/http"
)

type adaptor struct {
	RecordHandler *recordHandler.RecordHandler
}

func New(recordHandler *recordHandler.RecordHandler) *adaptor {
	return &adaptor{
		RecordHandler: recordHandler,
	}
}

type RetrieveRequest struct {
	WrappedIdentifier wrappedIdentifier.Wrapped `json:"identifier"`
}

type RetrieveResponse struct {
	ZX303StatusReading status.Reading `json:"zx303StatusReading"`
}

func (s *adaptor) Retrieve(r *http.Request, request *RetrieveRequest, response *RetrieveResponse) error {
	claims, err := wrappedClaims.UnwrapClaimsFromContext(r)
	if err != nil {
		log.Warn(err.Error())
		return err
	}

	retrieveZX303Response, err := s.RecordHandler.Retrieve(
		&recordHandler.RetrieveRequest{
			Claims:     claims,
			Identifier: request.WrappedIdentifier.Identifier,
		})
	if err != nil {
		return err
	}

	response.ZX303StatusReading = retrieveZX303Response.ZX303StatusReading

	return nil
}

type CollectRequest struct {
	Criteria []wrappedCriterion.Wrapped `json:"criteria"`
	Query    query.Query                `json:"query"`
}

type CollectResponse struct {
	Records []status.Reading `json:"records"`
	Total   int              `json:"total"`
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

	collectZX303Response, err := s.RecordHandler.Collect(&recordHandler.CollectRequest{
		Claims:   claims,
		Criteria: criteria,
		Query:    request.Query,
	})
	if err != nil {
		return err
	}

	response.Records = collectZX303Response.Records
	response.Total = collectZX303Response.Total
	return nil
}
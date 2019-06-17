package jsonRpc

import (
	"github.com/iot-my-world/brain/log"
	"github.com/iot-my-world/brain/search/criterion"
	wrappedCriterion "github.com/iot-my-world/brain/search/criterion/wrapped"
	wrappedIdentifier "github.com/iot-my-world/brain/search/identifier/wrapped"
	"github.com/iot-my-world/brain/search/query"
	wrappedClaims "github.com/iot-my-world/brain/security/claims/wrapped"
	"github.com/iot-my-world/brain/tracker/zx303"
	zx303RecordHandler "github.com/iot-my-world/brain/tracker/zx303/recordHandler"
	"net/http"
)

type adaptor struct {
	RecordHandler *zx303RecordHandler.RecordHandler
}

func New(recordHandler *zx303RecordHandler.RecordHandler) *adaptor {
	return &adaptor{
		RecordHandler: recordHandler,
	}
}

type RetrieveRequest struct {
	WrappedIdentifier wrappedIdentifier.Wrapped `json:"identifier"`
}

type RetrieveResponse struct {
	ZX303 zx303.ZX303 `json:"zx303"`
}

func (s *adaptor) Retrieve(r *http.Request, request *RetrieveRequest, response *RetrieveResponse) error {
	claims, err := wrappedClaims.UnwrapClaimsFromContext(r)
	if err != nil {
		log.Warn(err.Error())
		return err
	}

	retrieveZX303Response, err := s.RecordHandler.Retrieve(
		&zx303RecordHandler.RetrieveRequest{
			Claims:     claims,
			Identifier: request.WrappedIdentifier.Identifier,
		})
	if err != nil {
		return err
	}

	response.ZX303 = retrieveZX303Response.ZX303

	return nil
}

type CollectRequest struct {
	Criteria []wrappedCriterion.Wrapped `json:"criteria"`
	Query    query.Query                `json:"query"`
}

type CollectResponse struct {
	Records []zx303.ZX303 `json:"records"`
	Total   int           `json:"total"`
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

	collectZX303Response, err := s.RecordHandler.Collect(&zx303RecordHandler.CollectRequest{
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

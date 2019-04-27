package jsonRpc

import (
	"gitlab.com/iotTracker/brain/log"
	"gitlab.com/iotTracker/brain/search/criterion"
	wrappedCriterion "gitlab.com/iotTracker/brain/search/criterion/wrapped"
	wrappedIdentifier "gitlab.com/iotTracker/brain/search/identifier/wrapped"
	"gitlab.com/iotTracker/brain/search/query"
	wrappedClaims "gitlab.com/iotTracker/brain/security/claims/wrapped"
	zx303GPSReading "gitlab.com/iotTracker/brain/tracker/zx303/reading/gps"
	zx303GPSReadingRecordHandler "gitlab.com/iotTracker/brain/tracker/zx303/reading/gps/recordHandler"
	"net/http"
)

type adaptor struct {
	RecordHandler *zx303GPSReadingRecordHandler.RecordHandler
}

func New(recordHandler *zx303GPSReadingRecordHandler.RecordHandler) *adaptor {
	return &adaptor{
		RecordHandler: recordHandler,
	}
}

type RetrieveRequest struct {
	WrappedIdentifier wrappedIdentifier.Wrapped `json:"identifier"`
}

type RetrieveResponse struct {
	ZX303GPSReading zx303GPSReading.Reading `json:"zx303GPSReading"`
}

func (s *adaptor) Retrieve(r *http.Request, request *RetrieveRequest, response *RetrieveResponse) error {
	claims, err := wrappedClaims.UnwrapClaimsFromContext(r)
	if err != nil {
		log.Warn(err.Error())
		return err
	}

	retrieveZX303Response, err := s.RecordHandler.Retrieve(
		&zx303GPSReadingRecordHandler.RetrieveRequest{
			Claims:     claims,
			Identifier: request.WrappedIdentifier.Identifier,
		})
	if err != nil {
		return err
	}

	response.ZX303GPSReading = retrieveZX303Response.ZX303GPSReading

	return nil
}

type CollectRequest struct {
	Criteria []wrappedCriterion.Wrapped `json:"criteria"`
	Query    query.Query                `json:"query"`
}

type CollectResponse struct {
	Records []zx303GPSReading.Reading `json:"records"`
	Total   int                       `json:"total"`
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

	collectZX303Response, err := s.RecordHandler.Collect(&zx303GPSReadingRecordHandler.CollectRequest{
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

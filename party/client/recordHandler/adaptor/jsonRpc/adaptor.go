package client

import (
	"gitlab.com/iotTracker/brain/log"
	"gitlab.com/iotTracker/brain/party/client"
	clientRecordHandler "gitlab.com/iotTracker/brain/party/client/recordHandler"
	"gitlab.com/iotTracker/brain/search/criterion"
	wrappedCriterion "gitlab.com/iotTracker/brain/search/criterion/wrapped"
	wrappedIdentifier "gitlab.com/iotTracker/brain/search/identifier/wrapped"
	"gitlab.com/iotTracker/brain/search/query"
	"gitlab.com/iotTracker/brain/security/wrappedClaims"
	"net/http"
)

type adaptor struct {
	RecordHandler clientRecordHandler.RecordHandler
}

func New(recordHandler clientRecordHandler.RecordHandler) *adaptor {
	return &adaptor{
		RecordHandler: recordHandler,
	}
}

type RetrieveRequest struct {
	Identifier wrappedIdentifier.Wrapped `json:"identifier"`
}

type RetrieveResponse struct {
	Client client.Client `json:"client" bson:"client"`
}

func (s *adaptor) Retrieve(r *http.Request, request *RetrieveRequest, response *RetrieveResponse) error {
	claims, err := wrappedClaims.UnwrapClaimsFromContext(r)
	if err != nil {
		log.Warn(err.Error())
		return err
	}

	id, err := request.Identifier.UnWrap()
	if err != nil {
		return err
	}

	retrieveClientResponse, err := s.RecordHandler.Retrieve(
		&clientRecordHandler.RetrieveRequest{
			Claims:     claims,
			Identifier: id,
		})
	if err != nil {
		return err
	}

	response.Client = retrieveClientResponse.Client

	return nil
}

type CollectRequest struct {
	Criteria []wrappedCriterion.Wrapped `json:"criteria"`
	Query    query.Query                `json:"query"`
}

type CollectResponse struct {
	Records []client.Client `json:"records"`
	Total   int             `json:"total"`
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

	collectClientResponse, err := s.RecordHandler.Collect(&clientRecordHandler.CollectRequest{
		Criteria: criteria,
		Query:    request.Query,
		Claims:   claims,
	})
	if err != nil {
		return err
	}

	response.Records = collectClientResponse.Records
	response.Total = collectClientResponse.Total
	return nil
}

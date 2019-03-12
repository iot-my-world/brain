package client

import (
	"gitlab.com/iotTracker/brain/api"
	"gitlab.com/iotTracker/brain/log"
	"gitlab.com/iotTracker/brain/party/client"
	clientRecordHandler "gitlab.com/iotTracker/brain/party/client/recordHandler"
	"gitlab.com/iotTracker/brain/search/criterion"
	"gitlab.com/iotTracker/brain/search/query"
	"gitlab.com/iotTracker/brain/search/wrappedCriterion"
	"gitlab.com/iotTracker/brain/search/wrappedIdentifier"
	"gitlab.com/iotTracker/brain/security/wrappedClaims"
	"gitlab.com/iotTracker/brain/validate/reasonInvalid"
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

type CreateRequest struct {
	Client client.Client `json:"client"`
}

type CreateResponse struct {
	Client client.Client `json:"client"`
}

func (s *adaptor) Create(r *http.Request, request *CreateRequest, response *CreateResponse) error {
	claims, err := wrappedClaims.UnwrapClaimsFromContext(r)
	if err != nil {
		log.Warn(err.Error())
		return err
	}

	createClientResponse := clientRecordHandler.CreateResponse{}

	if err := s.RecordHandler.Create(
		&clientRecordHandler.CreateRequest{
			Client: request.Client,
			Claims: claims,
		},
		&createClientResponse); err != nil {
		return err
	}

	response.Client = createClientResponse.Client

	return nil
}

type RetrieveRequest struct {
	Identifier wrappedIdentifier.WrappedIdentifier `json:"identifier"`
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

	retrieveClientResponse := clientRecordHandler.RetrieveResponse{}
	if err := s.RecordHandler.Retrieve(
		&clientRecordHandler.RetrieveRequest{
			Claims:     claims,
			Identifier: id,
		},
		&retrieveClientResponse); err != nil {
		return err
	}

	response.Client = retrieveClientResponse.Client

	return nil
}

type UpdateRequest struct {
	Identifier wrappedIdentifier.WrappedIdentifier `json:"identifier"`
	Client     client.Client                       `json:"client"`
}

type UpdateResponse struct {
	Client client.Client `json:"client"`
}

func (s *adaptor) Update(r *http.Request, request *UpdateRequest, response *UpdateResponse) error {
	claims, err := wrappedClaims.UnwrapClaimsFromContext(r)
	if err != nil {
		log.Warn(err.Error())
		return err
	}

	id, err := request.Identifier.UnWrap()
	if err != nil {
		return err
	}

	updateClientResponse := clientRecordHandler.UpdateResponse{}
	if err := s.RecordHandler.Update(
		&clientRecordHandler.UpdateRequest{
			Claims:     claims,
			Identifier: id,
		},
		&updateClientResponse); err != nil {
		return err
	}

	response.Client = updateClientResponse.Client

	return nil
}

type DeleteRequest struct {
	Identifier wrappedIdentifier.WrappedIdentifier `json:"identifier"`
}

type DeleteResponse struct {
	Client client.Client `json:"client"`
}

func (s *adaptor) Delete(r *http.Request, request *DeleteRequest, response *DeleteResponse) error {
	claims, err := wrappedClaims.UnwrapClaimsFromContext(r)
	if err != nil {
		log.Warn(err.Error())
		return err
	}

	id, err := request.Identifier.UnWrap()
	if err != nil {
		return err
	}

	deleteClientResponse := clientRecordHandler.DeleteResponse{}
	if err := s.RecordHandler.Delete(
		&clientRecordHandler.DeleteRequest{
			Claims:     claims,
			Identifier: id,
		},
		&deleteClientResponse); err != nil {
		return err
	}

	response.Client = deleteClientResponse.Client

	return nil
}

type ValidateRequest struct {
	Client client.Client `json:"client"`
	Method api.Method    `json:"method"`
}

type ValidateResponse struct {
	ReasonsInvalid []reasonInvalid.ReasonInvalid `json:"reasonsInvalid"`
}

func (s *adaptor) Validate(r *http.Request, request *ValidateRequest, response *ValidateResponse) error {
	claims, err := wrappedClaims.UnwrapClaimsFromContext(r)
	if err != nil {
		log.Warn(err.Error())
		return err
	}

	validateClientResponse := clientRecordHandler.ValidateResponse{}
	if err := s.RecordHandler.Validate(
		&clientRecordHandler.ValidateRequest{
			Claims: claims,
			Client: request.Client,
			Method: request.Method,
		},
		&validateClientResponse); err != nil {
		return err
	}

	response.ReasonsInvalid = validateClientResponse.ReasonsInvalid

	return nil
}

type CollectRequest struct {
	Criteria []wrappedCriterion.WrappedCriterion `json:"criteria"`
	Query    query.Query                         `json:"query"`
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

	collectClientResponse := clientRecordHandler.CollectResponse{}
	if err := s.RecordHandler.Collect(&clientRecordHandler.CollectRequest{
		Criteria: criteria,
		Query:    request.Query,
		Claims:   claims,
	},
		&collectClientResponse); err != nil {
		return err
	}

	response.Records = collectClientResponse.Records
	response.Total = collectClientResponse.Total
	return nil
}

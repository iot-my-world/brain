package jsonRpc

import (
	"gitlab.com/iotTracker/brain/api"
	"gitlab.com/iotTracker/brain/log"
	"gitlab.com/iotTracker/brain/search/criterion"
	"gitlab.com/iotTracker/brain/search/query"
	"gitlab.com/iotTracker/brain/search/wrappedCriterion"
	"gitlab.com/iotTracker/brain/search/wrappedIdentifier"
	"gitlab.com/iotTracker/brain/security/wrappedClaims"
	"gitlab.com/iotTracker/brain/tracker/device/tk102"
	tk102RecordHandler "gitlab.com/iotTracker/brain/tracker/device/tk102/recordHandler"
	"gitlab.com/iotTracker/brain/validate/reasonInvalid"
	"net/http"
)

type adaptor struct {
	RecordHandler tk102RecordHandler.RecordHandler
}

func New(recordHandler tk102RecordHandler.RecordHandler) *adaptor {
	return &adaptor{
		RecordHandler: recordHandler,
	}
}

type CreateRequest struct {
	TK102 tk102.TK102 `json:"tk102"`
}

type CreateResponse struct {
	TK102 tk102.TK102 `json:"tk102"`
}

func (s *adaptor) Create(r *http.Request, request *CreateRequest, response *CreateResponse) error {
	claims, err := wrappedClaims.UnwrapClaimsFromContext(r)
	if err != nil {
		log.Warn(err.Error())
		return err
	}

	createTK102Response := tk102RecordHandler.CreateResponse{}

	if err := s.RecordHandler.Create(
		&tk102RecordHandler.CreateRequest{
			TK102:  request.TK102,
			Claims: claims,
		},
		&createTK102Response); err != nil {
		return err
	}

	response.TK102 = createTK102Response.TK102

	return nil
}

type RetrieveRequest struct {
	Identifier wrappedIdentifier.WrappedIdentifier `json:"identifier"`
}

type RetrieveResponse struct {
	TK102 tk102.TK102 `json:"tk102"`
}

func (s *adaptor) Retrieve(r *http.Request, request *RetrieveRequest, response *RetrieveResponse) error {
	id, err := request.Identifier.UnWrap()
	if err != nil {
		return err
	}

	retrieveTK102Response := tk102RecordHandler.RetrieveResponse{}
	if err := s.RecordHandler.Retrieve(
		&tk102RecordHandler.RetrieveRequest{
			Identifier: id,
		},
		&retrieveTK102Response); err != nil {
		return err
	}

	response.TK102 = retrieveTK102Response.TK102

	return nil
}

type DeleteRequest struct {
	Identifier wrappedIdentifier.WrappedIdentifier `json:"identifier"`
}

type DeleteResponse struct {
	TK102 tk102.TK102 `json:"tk102"`
}

func (s *adaptor) Delete(r *http.Request, request *DeleteRequest, response *DeleteResponse) error {
	id, err := request.Identifier.UnWrap()
	if err != nil {
		return err
	}

	deleteTK102Response := tk102RecordHandler.DeleteResponse{}
	if err := s.RecordHandler.Delete(
		&tk102RecordHandler.DeleteRequest{
			Identifier: id,
		},
		&deleteTK102Response); err != nil {
		return err
	}

	response.TK102 = deleteTK102Response.TK102

	return nil
}

type ValidateRequest struct {
	TK102  tk102.TK102 `json:"tk102"`
	Method api.Method  `json:"method"`
}

type ValidateResponse struct {
	ReasonsInvalid []reasonInvalid.ReasonInvalid `json:"reasonsInvalid"`
}

func (s *adaptor) Validate(r *http.Request, request *ValidateRequest, response *ValidateResponse) error {

	validateTK102Response := tk102RecordHandler.ValidateResponse{}
	if err := s.RecordHandler.Validate(
		&tk102RecordHandler.ValidateRequest{
			TK102:  request.TK102,
			Method: request.Method,
		},
		&validateTK102Response); err != nil {
		return err
	}

	response.ReasonsInvalid = validateTK102Response.ReasonsInvalid

	return nil
}

type CollectRequest struct {
	Criteria []wrappedCriterion.WrappedCriterion `json:"criteria"`
	Query    query.Query                         `json:"query"`
}

type CollectResponse struct {
	Records []tk102.TK102 `json:"records"`
	Total   int           `json:"total"`
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

	collectTK102Response := tk102RecordHandler.CollectResponse{}
	if err := s.RecordHandler.Collect(&tk102RecordHandler.CollectRequest{
		Criteria: criteria,
		Query:    request.Query,
	},
		&collectTK102Response); err != nil {
		return err
	}

	response.Records = collectTK102Response.Records
	response.Total = collectTK102Response.Total
	return nil
}

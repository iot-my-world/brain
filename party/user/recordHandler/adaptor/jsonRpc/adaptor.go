package user

import (
	"gitlab.com/iotTracker/brain/log"
	"gitlab.com/iotTracker/brain/party/user"
	userRecordHandler "gitlab.com/iotTracker/brain/party/user/recordHandler"
	"gitlab.com/iotTracker/brain/search/criterion"
	"gitlab.com/iotTracker/brain/search/query"
	"gitlab.com/iotTracker/brain/search/wrappedCriterion"
	"gitlab.com/iotTracker/brain/search/wrappedIdentifier"
	"gitlab.com/iotTracker/brain/security/wrappedClaims"
	"net/http"
)

type adaptor struct {
	RecordHandler userRecordHandler.RecordHandler
}

func New(recordHandler userRecordHandler.RecordHandler) *adaptor {
	return &adaptor{
		RecordHandler: recordHandler,
	}
}

type CreateRequest struct {
	User user.User `json:"user"`
}

type CreateResponse struct {
	User user.User `json:"user"`
}

func (s *adaptor) Create(r *http.Request, request *CreateRequest, response *CreateResponse) error {
	claims, err := wrappedClaims.UnwrapClaimsFromContext(r)
	if err != nil {
		log.Warn(err.Error())
		return err
	}

	createUserResponse := userRecordHandler.CreateResponse{}
	if err := s.RecordHandler.Create(
		&userRecordHandler.CreateRequest{
			Claims: claims,
			User:   request.User,
		}, &createUserResponse); err != nil {
		return err
	}

	response.User = createUserResponse.User

	return nil
}

type RetrieveRequest struct {
	Identifier wrappedIdentifier.WrappedIdentifier `json:"identifier"`
}

type RetrieveResponse struct {
	User user.User `json:"user" bson:"user"`
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

	retrieveUserResponse := userRecordHandler.RetrieveResponse{}
	if err := s.RecordHandler.Retrieve(
		&userRecordHandler.RetrieveRequest{
			Claims:     claims,
			Identifier: id,
		},
		&retrieveUserResponse); err != nil {
		return err
	}

	response.User = retrieveUserResponse.User

	return nil
}

type CollectRequest struct {
	Criteria []wrappedCriterion.WrappedCriterion `json:"criteria"`
	Query    query.Query                         `json:"query"`
}

type CollectResponse struct {
	Records []user.User `json:"records"`
	Total   int         `json:"total"`
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

	collectCompanyResponse := userRecordHandler.CollectResponse{}
	if err := s.RecordHandler.Collect(&userRecordHandler.CollectRequest{
		Criteria: criteria,
		Query:    request.Query,
		Claims:   claims,
	},
		&collectCompanyResponse); err != nil {
		return err
	}

	response.Records = collectCompanyResponse.Records
	response.Total = collectCompanyResponse.Total
	return nil
}

package user

import (
	"github.com/iot-my-world/brain/log"
	"github.com/iot-my-world/brain/search/criterion"
	wrappedCriterion "github.com/iot-my-world/brain/search/criterion/wrapped"
	wrappedIdentifier "github.com/iot-my-world/brain/search/identifier/wrapped"
	"github.com/iot-my-world/brain/search/query"
	wrappedClaims "github.com/iot-my-world/brain/security/claims/wrapped"
	humanUser "github.com/iot-my-world/brain/user/human"
	userRecordHandler "github.com/iot-my-world/brain/user/human/recordHandler"
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

type RetrieveRequest struct {
	WrappedIdentifier wrappedIdentifier.Wrapped `json:"identifier"`
}

type RetrieveResponse struct {
	User humanUser.User `json:"user" bson:"user"`
}

func (s *adaptor) Retrieve(r *http.Request, request *RetrieveRequest, response *RetrieveResponse) error {
	claims, err := wrappedClaims.UnwrapClaimsFromContext(r)
	if err != nil {
		log.Warn(err.Error())
		return err
	}

	retrieveUserResponse, err := s.RecordHandler.Retrieve(
		&userRecordHandler.RetrieveRequest{
			Claims:     claims,
			Identifier: request.WrappedIdentifier.Identifier,
		})
	if err != nil {
		return err
	}

	response.User = retrieveUserResponse.User

	return nil
}

type CollectRequest struct {
	Criteria []wrappedCriterion.Wrapped `json:"criteria"`
	Query    query.Query                `json:"query"`
}

type CollectResponse struct {
	Records []humanUser.User `json:"records"`
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

	collectCompanyResponse, err := s.RecordHandler.Collect(&userRecordHandler.CollectRequest{
		Criteria: criteria,
		Query:    request.Query,
		Claims:   claims,
	})
	if err != nil {
		return err
	}

	response.Records = collectCompanyResponse.Records
	response.Total = collectCompanyResponse.Total
	return nil
}

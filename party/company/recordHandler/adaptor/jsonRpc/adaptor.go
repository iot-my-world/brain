package company

import (
	"gitlab.com/iotTracker/brain/log"
	"gitlab.com/iotTracker/brain/party/company"
	companyRecordHandler "gitlab.com/iotTracker/brain/party/company/recordHandler"
	"gitlab.com/iotTracker/brain/search/criterion"
	wrappedCriterion "gitlab.com/iotTracker/brain/search/criterion/wrapped"
	wrappedIdentifier "gitlab.com/iotTracker/brain/search/identifier/wrapped"
	"gitlab.com/iotTracker/brain/search/query"
	"gitlab.com/iotTracker/brain/security/wrappedClaims"
	"net/http"
)

type adaptor struct {
	RecordHandler companyRecordHandler.RecordHandler
}

func New(recordHandler companyRecordHandler.RecordHandler) *adaptor {
	return &adaptor{
		RecordHandler: recordHandler,
	}
}

type RetrieveRequest struct {
	Identifier wrappedIdentifier.Wrapped `json:"identifier"`
}

type RetrieveResponse struct {
	Company company.Company `json:"company"`
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

	retrieveCompanyResponse, err := s.RecordHandler.Retrieve(
		&companyRecordHandler.RetrieveRequest{
			Claims:     claims,
			Identifier: id,
		})
	if err != nil {
		return err
	}

	response.Company = retrieveCompanyResponse.Company

	return nil
}

type CollectRequest struct {
	Criteria []wrappedCriterion.Wrapped `json:"criteria"`
	Query    query.Query                `json:"query"`
}

type CollectResponse struct {
	Records []company.Company `json:"records"`
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

	collectCompanyResponse, err := s.RecordHandler.Collect(&companyRecordHandler.CollectRequest{
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

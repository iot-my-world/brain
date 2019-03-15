package company

import (
	"gitlab.com/iotTracker/brain/api"
	"gitlab.com/iotTracker/brain/log"
	"gitlab.com/iotTracker/brain/party/company"
	companyRecordHandler "gitlab.com/iotTracker/brain/party/company/recordHandler"
	"gitlab.com/iotTracker/brain/search/criterion"
	"gitlab.com/iotTracker/brain/search/query"
	"gitlab.com/iotTracker/brain/search/wrappedCriterion"
	"gitlab.com/iotTracker/brain/search/wrappedIdentifier"
	"gitlab.com/iotTracker/brain/security/wrappedClaims"
	"gitlab.com/iotTracker/brain/validate/reasonInvalid"
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

type CreateRequest struct {
	Company company.Company `json:"company"`
}

type CreateResponse struct {
	Company company.Company `json:"company"`
}

func (s *adaptor) Create(r *http.Request, request *CreateRequest, response *CreateResponse) error {
	claims, err := wrappedClaims.UnwrapClaimsFromContext(r)
	if err != nil {
		log.Warn(err.Error())
		return err
	}

	createCompanyResponse := companyRecordHandler.CreateResponse{}

	if err := s.RecordHandler.Create(
		&companyRecordHandler.CreateRequest{
			Company: request.Company,
			Claims:  claims,
		},
		&createCompanyResponse); err != nil {
		return err
	}

	response.Company = createCompanyResponse.Company

	return nil
}

type RetrieveRequest struct {
	Identifier wrappedIdentifier.WrappedIdentifier `json:"identifier"`
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

	retrieveCompanyResponse := companyRecordHandler.RetrieveResponse{}
	if err := s.RecordHandler.Retrieve(
		&companyRecordHandler.RetrieveRequest{
			Claims:     claims,
			Identifier: id,
		},
		&retrieveCompanyResponse); err != nil {
		return err
	}

	response.Company = retrieveCompanyResponse.Company

	return nil
}

type ValidateRequest struct {
	Company company.Company `json:"company"`
	Method  api.Method      `json:"method"`
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

	validateCompanyResponse := companyRecordHandler.ValidateResponse{}
	if err := s.RecordHandler.Validate(&companyRecordHandler.ValidateRequest{
		Claims:  claims,
		Company: request.Company,
		Method:  request.Method,
	},
		&validateCompanyResponse); err != nil {
		return err
	}

	response.ReasonsInvalid = validateCompanyResponse.ReasonsInvalid

	return nil
}

type CollectRequest struct {
	Criteria []wrappedCriterion.WrappedCriterion `json:"criteria"`
	Query    query.Query                         `json:"query"`
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

	collectCompanyResponse := companyRecordHandler.CollectResponse{}
	if err := s.RecordHandler.Collect(&companyRecordHandler.CollectRequest{
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

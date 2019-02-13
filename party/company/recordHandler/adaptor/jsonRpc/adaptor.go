package company

import (
	"gitlab.com/iotTracker/brain/api"
	"gitlab.com/iotTracker/brain/party/company"
	companyRecordHandler "gitlab.com/iotTracker/brain/party/company/recordHandler"
	"gitlab.com/iotTracker/brain/search/criterion"
	"gitlab.com/iotTracker/brain/search/query"
	"gitlab.com/iotTracker/brain/search/wrappedCriterion"
	"gitlab.com/iotTracker/brain/search/wrappedIdentifier"
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
	createCompanyResponse := companyRecordHandler.CreateResponse{}

	if err := s.RecordHandler.Create(
		&companyRecordHandler.CreateRequest{
			Company: request.Company,
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
	Company company.Company `json:"company" bson:"company"`
}

func (s *adaptor) Retrieve(r *http.Request, request *RetrieveRequest, response *RetrieveResponse) error {
	id, err := request.Identifier.UnWrap()
	if err != nil {
		return err
	}

	retrieveCompanyResponse := companyRecordHandler.RetrieveResponse{}
	if err := s.RecordHandler.Retrieve(
		&companyRecordHandler.RetrieveRequest{
			Identifier: id,
		},
		&retrieveCompanyResponse); err != nil {
		return err
	}

	response.Company = retrieveCompanyResponse.Company

	return nil
}

type UpdateRequest struct {
	Identifier wrappedIdentifier.WrappedIdentifier `json:"identifier"`
	Company    company.Company                     `json:"company"`
}

type UpdateResponse struct {
	Company company.Company `json:"company"`
}

func (s *adaptor) Update(r *http.Request, request *UpdateRequest, response *UpdateResponse) error {
	id, err := request.Identifier.UnWrap()
	if err != nil {
		return err
	}

	updateCompanyResponse := companyRecordHandler.UpdateResponse{}
	if err := s.RecordHandler.Update(
		&companyRecordHandler.UpdateRequest{
			Identifier: id,
		},
		&updateCompanyResponse); err != nil {
		return err
	}

	response.Company = updateCompanyResponse.Company

	return nil
}

type DeleteRequest struct {
	Identifier wrappedIdentifier.WrappedIdentifier `json:"identifier"`
}

type DeleteResponse struct {
	Company company.Company `json:"company"`
}

func (s *adaptor) Delete(r *http.Request, request *DeleteRequest, response *DeleteResponse) error {
	id, err := request.Identifier.UnWrap()
	if err != nil {
		return err
	}

	deleteCompanyResponse := companyRecordHandler.DeleteResponse{}
	if err := s.RecordHandler.Delete(
		&companyRecordHandler.DeleteRequest{
			Identifier: id,
		},
		&deleteCompanyResponse); err != nil {
		return err
	}

	response.Company = deleteCompanyResponse.Company

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

	validateCompanyResponse := companyRecordHandler.ValidateResponse{}
	if err := s.RecordHandler.Validate(
		&companyRecordHandler.ValidateRequest{
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
	// unwrap criteria
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
	},
		&collectCompanyResponse); err != nil {
		return err
	}

	response.Records = collectCompanyResponse.Records
	response.Total = collectCompanyResponse.Total
	return nil
}

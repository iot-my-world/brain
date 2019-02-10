package company

import (
	"gitlab.com/iotTracker/brain/api"
	"gitlab.com/iotTracker/brain/party/company"
	"gitlab.com/iotTracker/brain/party/company/recordHandler"
	"gitlab.com/iotTracker/brain/search"
	"gitlab.com/iotTracker/brain/validate/reasonInvalid"
	"net/http"
)

type adaptor struct {
	RecordHandler recordHandler.RecordHandler
}

func New(recordHandler recordHandler.RecordHandler) *adaptor {
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
	createCompanyResponse := recordHandler.CreateResponse{}

	if err := s.RecordHandler.Create(
		&recordHandler.CreateRequest{
			Company: request.Company,
		},
		&createCompanyResponse); err != nil {
		return err
	}

	response.Company = createCompanyResponse.Company

	return nil
}

type RetrieveRequest struct {
	Identifier search.WrappedIdentifier `json:"identifier"`
}

type RetrieveResponse struct {
	Company company.Company `json:"company" bson:"company"`
}

func (s *adaptor) Retrieve(r *http.Request, request *RetrieveRequest, response *RetrieveResponse) error {
	id, err := request.Identifier.UnWrap()
	if err != nil {
		return err
	}

	retrieveCompanyResponse := recordHandler.RetrieveResponse{}
	if err := s.RecordHandler.Retrieve(
		&recordHandler.RetrieveRequest{
			Identifier: id,
		},
		&retrieveCompanyResponse); err != nil {
		return err
	}

	response.Company = retrieveCompanyResponse.Company

	return nil
}

type UpdateRequest struct {
	Identifier search.WrappedIdentifier `json:"identifier"`
	Company    company.Company          `json:"company"`
}

type UpdateResponse struct {
	Company company.Company `json:"company"`
}

func (s *adaptor) Update(r *http.Request, request *UpdateRequest, response *UpdateResponse) error {
	id, err := request.Identifier.UnWrap()
	if err != nil {
		return err
	}

	updateCompanyResponse := recordHandler.UpdateResponse{}
	if err := s.RecordHandler.Update(
		&recordHandler.UpdateRequest{
			Identifier: id,
		},
		&updateCompanyResponse); err != nil {
		return err
	}

	response.Company = updateCompanyResponse.Company

	return nil
}

type DeleteRequest struct {
	Identifier search.WrappedIdentifier `json:"identifier"`
}

type DeleteResponse struct {
	Company company.Company `json:"company"`
}

func (s *adaptor) Delete(r *http.Request, request *DeleteRequest, response *DeleteResponse) error {
	id, err := request.Identifier.UnWrap()
	if err != nil {
		return err
	}

	deleteCompanyResponse := recordHandler.DeleteResponse{}
	if err := s.RecordHandler.Delete(
		&recordHandler.DeleteRequest{
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

	validateCompanyResponse := recordHandler.ValidateResponse{}
	if err := s.RecordHandler.Validate(
		&recordHandler.ValidateRequest{
			Company: request.Company,
			Method:  request.Method,
		},
		&validateCompanyResponse); err != nil {
		return err
	}

	response.ReasonsInvalid = validateCompanyResponse.ReasonsInvalid

	return nil
}

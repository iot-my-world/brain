package company

import (
	"gitlab.com/iotTracker/brain/party/company"
	"gitlab.com/iotTracker/brain/search"
	"gitlab.com/iotTracker/brain/validate/reasonInvalid"
	"net/http"
)

type adaptor struct {
	RecordHandler company.RecordHandler
}

func New(recordHandler company.RecordHandler) *adaptor {
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
	createCompanyResponse := company.CreateResponse{}

	if err := s.RecordHandler.Create(
		&company.CreateRequest{
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

	retrieveCompanyResponse := company.RetrieveResponse{}
	if err := s.RecordHandler.Retrieve(
		&company.RetrieveRequest{
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

	updateCompanyResponse := company.UpdateResponse{}
	if err := s.RecordHandler.Update(
		&company.UpdateRequest{
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

	deleteCompanyResponse := company.DeleteResponse{}
	if err := s.RecordHandler.Delete(
		&company.DeleteRequest{
			Identifier: id,
		},
		&deleteCompanyResponse); err != nil {
		return err
	}

	response.Company = deleteCompanyResponse.Company

	return nil
}

type ValidateRequest struct {
	Company             company.Company `json:"company"`
	IgnoreReasonsMethod string          `json:"ignoreReasonsMethod"`
}

type ValidateResponse struct {
	ReasonsInvalid []reasonInvalid.ReasonInvalid `json:"reasonsInvalid"`
}

func (s *adaptor) Validate(r *http.Request, request *ValidateRequest, response *ValidateResponse) error {

	validateCompanyResponse := company.ValidateResponse{}
	if err := s.RecordHandler.Validate(
		&company.ValidateRequest{
			Company:             request.Company,
			IgnoreReasonsMethod: request.IgnoreReasonsMethod,
		},
		&validateCompanyResponse); err != nil {
		return err
	}

	response.ReasonsInvalid = validateCompanyResponse.ReasonsInvalid

	return nil
}

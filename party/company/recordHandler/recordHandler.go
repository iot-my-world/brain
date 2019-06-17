package recordHandler

import (
	"gitlab.com/iotTracker/brain/party/company"
	companyRecordHandlerException "gitlab.com/iotTracker/brain/party/company/recordHandler/exception"
	brainRecordHandler "gitlab.com/iotTracker/brain/recordHandler"
	brainRecordHandlerException "gitlab.com/iotTracker/brain/recordHandler/exception"
	"gitlab.com/iotTracker/brain/search/criterion"
	"gitlab.com/iotTracker/brain/search/identifier"
	"gitlab.com/iotTracker/brain/search/query"
	"gitlab.com/iotTracker/brain/security/claims"
)

type RecordHandler struct {
	recordHandler brainRecordHandler.RecordHandler
}

func New(
	brainCompanyRecordHandler brainRecordHandler.RecordHandler,
) *RecordHandler {

	return &RecordHandler{
		recordHandler: brainCompanyRecordHandler,
	}
}

type CreateRequest struct {
	Company company.Company
}

type CreateResponse struct {
	Company company.Company
}

func (r *RecordHandler) Create(request *CreateRequest) (*CreateResponse, error) {
	createResponse := brainRecordHandler.CreateResponse{}
	if err := r.recordHandler.Create(&brainRecordHandler.CreateRequest{
		Entity: &request.Company,
	}, &createResponse); err != nil {
		return nil, companyRecordHandlerException.Create{Reasons: []string{err.Error()}}
	}
	createdCompany, ok := createResponse.Entity.(*company.Company)
	if !ok {
		return nil, companyRecordHandlerException.Create{Reasons: []string{"could not cast created entity to company"}}
	}

	return &CreateResponse{
		Company: *createdCompany,
	}, nil
}

type RetrieveRequest struct {
	Claims     claims.Claims
	Identifier identifier.Identifier
}

type RetrieveResponse struct {
	Company company.Company
}

func (r *RecordHandler) Retrieve(request *RetrieveRequest) (*RetrieveResponse, error) {
	retrievedCompany := company.Company{}
	retrieveResponse := brainRecordHandler.RetrieveResponse{
		Entity: &retrievedCompany,
	}
	if err := r.recordHandler.Retrieve(&brainRecordHandler.RetrieveRequest{
		Claims:     request.Claims,
		Identifier: request.Identifier,
	}, &retrieveResponse); err != nil {
		switch err.(type) {
		case brainRecordHandlerException.NotFound:
			return nil, companyRecordHandlerException.NotFound{}
		default:
			return nil, err
		}
	}

	return &RetrieveResponse{
		Company: retrievedCompany,
	}, nil
}

type UpdateRequest struct {
	Claims     claims.Claims
	Identifier identifier.Identifier
	Company    company.Company
}

type UpdateResponse struct{}

func (r *RecordHandler) Update(request *UpdateRequest) (*UpdateResponse, error) {
	updateResponse := brainRecordHandler.UpdateResponse{}
	if err := r.recordHandler.Update(&brainRecordHandler.UpdateRequest{
		Claims:     request.Claims,
		Identifier: request.Identifier,
		Entity:     &request.Company,
	}, &updateResponse); err != nil {
		return nil, companyRecordHandlerException.Update{Reasons: []string{err.Error()}}
	}

	return &UpdateResponse{}, nil
}

type DeleteRequest struct {
	Claims     claims.Claims
	Identifier identifier.Identifier
}

type DeleteResponse struct {
}

func (r *RecordHandler) Delete(request *DeleteRequest) (*DeleteResponse, error) {
	deleteResponse := brainRecordHandler.DeleteResponse{}
	if err := r.recordHandler.Delete(&brainRecordHandler.DeleteRequest{
		Claims:     request.Claims,
		Identifier: request.Identifier,
	}, &deleteResponse); err != nil {
		return nil, companyRecordHandlerException.Delete{Reasons: []string{err.Error()}}
	}

	return &DeleteResponse{}, nil
}

type CollectRequest struct {
	Claims   claims.Claims
	Criteria []criterion.Criterion
	Query    query.Query
}

type CollectResponse struct {
	Records []company.Company
	Total   int
}

func (r *RecordHandler) Collect(request *CollectRequest) (*CollectResponse, error) {
	var collectedCompanies []company.Company
	collectResponse := brainRecordHandler.CollectResponse{
		Records: &collectedCompanies,
	}
	err := r.recordHandler.Collect(&brainRecordHandler.CollectRequest{
		Claims:   request.Claims,
		Criteria: request.Criteria,
		Query:    request.Query,
	}, &collectResponse)
	if err != nil {
		return nil, companyRecordHandlerException.Collect{Reasons: []string{err.Error()}}
	}

	return &CollectResponse{
		Records: collectedCompanies,
		Total:   collectResponse.Total,
	}, nil
}

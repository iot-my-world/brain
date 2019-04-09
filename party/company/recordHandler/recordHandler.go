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
) (*RecordHandler, error) {

	if brainCompanyRecordHandler == nil {
		return nil, companyRecordHandlerException.RecordHandlerNil{}
	}
	return &RecordHandler{
		recordHandler: brainCompanyRecordHandler,
	}, nil
}

type CreateRequest struct {
	Company company.Company
}

type CreateResponse struct {
	Company company.Company
}

func (r *RecordHandler) Create(request *CreateRequest) (*CreateResponse, error) {
	createResponse, err := r.recordHandler.Create(&brainRecordHandler.CreateRequest{
		Entity: request.Company,
	})
	if err != nil {
		return nil, companyRecordHandlerException.Create{Reasons: []string{err.Error()}}
	}

	createdCompany, ok := createResponse.Entity.(company.Company)
	if !ok {
		return nil, companyRecordHandlerException.Create{Reasons: []string{"could not cast created entity to company"}}
	}

	return &CreateResponse{Company: createdCompany}, nil
}

type RetrieveRequest struct {
	Claims     claims.Claims
	Identifier identifier.Identifier
}

type RetrieveResponse struct {
	Company company.Company
}

func (r *RecordHandler) Retrieve(request *RetrieveRequest) (*RetrieveResponse, error) {
	retrieveResponse, err := r.recordHandler.Retrieve(&brainRecordHandler.RetrieveRequest{
		Claims:     request.Claims,
		Identifier: request.Identifier,
	})
	if err != nil {
		switch err.(type) {
		case brainRecordHandlerException.NotFound:
			return nil, companyRecordHandlerException.NotFound{}
		default:
			return nil, err
		}
	}

	retrievedCompany, ok := retrieveResponse.Entity.(company.Company)
	if !ok {
		return nil, companyRecordHandlerException.Retrieval{Reasons: []string{"could not case retrieved entity to company"}}
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

type UpdateResponse struct {
	Company company.Company
}

type DeleteRequest struct {
	Claims     claims.Claims
	Identifier identifier.Identifier
}

type DeleteResponse struct {
	Company company.Company
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

package recordHandler

import (
	"gitlab.com/iotTracker/brain/log"
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

	if brainCompanyRecordHandler == nil {
		log.Fatal(companyRecordHandlerException.RecordHandlerNil{}.Error())
	}
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
	createResponse, err := r.recordHandler.Create(&brainRecordHandler.CreateRequest{
		Entity: request.Company,
	})
	if err != nil {
		return nil, companyRecordHandlerException.Create{Reasons: []string{err.Error()}}
	}

	if createResponse.Entity == nil {
		return nil, companyRecordHandlerException.Create{Reasons: []string{"created entity is nil"}}
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

	if retrieveResponse.Entity == nil {
		return nil, companyRecordHandlerException.Retrieve{Reasons: []string{"retrieved entity is nil"}}
	}
	retrievedCompany, ok := retrieveResponse.Entity.(company.Company)
	if !ok {
		return nil, companyRecordHandlerException.Retrieve{Reasons: []string{"could not case retrieved entity to company"}}
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

func (r *RecordHandler) Update(request *UpdateRequest) (*UpdateResponse, error) {
	updateResponse, err := r.recordHandler.Update(&brainRecordHandler.UpdateRequest{
		Claims:     request.Claims,
		Identifier: request.Identifier,
		Entity:     request.Company,
	})
	if err != nil {
		return nil, companyRecordHandlerException.Update{Reasons: []string{err.Error()}}
	}

	if updateResponse.Entity == nil {
		return nil, companyRecordHandlerException.Update{Reasons: []string{"updated entity is nil"}}
	}
	updatedCompany, ok := updateResponse.Entity.(company.Company)
	if !ok {
		return nil, companyRecordHandlerException.Update{Reasons: []string{"could not cast updated entity to company"}}
	}

	return &UpdateResponse{Company: updatedCompany}, nil
}

type DeleteRequest struct {
	Claims     claims.Claims
	Identifier identifier.Identifier
}

type DeleteResponse struct {
	Company company.Company
}

func (r *RecordHandler) Delete(request *DeleteRequest) (*DeleteResponse, error) {
	deleteResponse, err := r.recordHandler.Delete(&brainRecordHandler.DeleteRequest{
		Claims:     request.Claims,
		Identifier: request.Identifier,
	})
	if err != nil {
		return nil, companyRecordHandlerException.Delete{Reasons: []string{err.Error()}}
	}

	if deleteResponse.Entity == nil {
		return nil, companyRecordHandlerException.Delete{Reasons: []string{"updated entity is nil"}}
	}
	deletedCompany, ok := deleteResponse.Entity.(company.Company)
	if !ok {
		return nil, companyRecordHandlerException.Delete{Reasons: []string{"could not cast deleted entity to company"}}
	}

	return &DeleteResponse{Company: deletedCompany}, nil
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
	collectResponse, err := r.recordHandler.Collect(&brainRecordHandler.CollectRequest{
		Claims:   request.Claims,
		Criteria: request.Criteria,
		Query:    request.Query,
	})
	if err != nil {
		return nil, companyRecordHandlerException.Delete{Reasons: []string{err.Error()}}
	}

	collectedCompanies := make([]company.Company, 0)
	if collectResponse.Records == nil {
		return nil, companyRecordHandlerException.Collect{Reasons: []string{"entities are nil in collect response"}}
	} else {
		for _, companyEntity := range collectResponse.Records {

			if companyEntity == nil {
				return nil, companyRecordHandlerException.Collect{Reasons: []string{"a collected entity is nil"}}
			}
			collectedCompany, ok := companyEntity.(company.Company)
			if !ok {
				return nil, companyRecordHandlerException.Collect{Reasons: []string{"could not cast a collected entity to company"}}
			}
			collectedCompanies = append(collectedCompanies, collectedCompany)
		}
	}

	return &CollectResponse{
		Records: collectedCompanies,
		Total:   collectResponse.Total,
	}, nil
}

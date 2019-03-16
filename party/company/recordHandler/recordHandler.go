package recordHandler

import (
	"gitlab.com/iotTracker/brain/party/company"
	"gitlab.com/iotTracker/brain/search/criterion"
	"gitlab.com/iotTracker/brain/search/identifier"
	"gitlab.com/iotTracker/brain/search/query"
	"gitlab.com/iotTracker/brain/security/claims"
)

type RecordHandler interface {
	Create(request *CreateRequest, response *CreateResponse) error
	Retrieve(request *RetrieveRequest, response *RetrieveResponse) error
	Update(request *UpdateRequest, response *UpdateResponse) error
	Delete(request *DeleteRequest, response *DeleteResponse) error
	Collect(request *CollectRequest, response *CollectResponse) error
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

type CreateRequest struct {
	Claims  claims.Claims
	Company company.Company
}

type CreateResponse struct {
	Company company.Company
}

type DeleteRequest struct {
	Claims     claims.Claims
	Identifier identifier.Identifier
}

type DeleteResponse struct {
	Company company.Company
}

type UpdateRequest struct {
	Claims     claims.Claims
	Identifier identifier.Identifier
	Company    company.Company
}

type UpdateResponse struct {
	Company company.Company
}

type RetrieveRequest struct {
	Claims     claims.Claims
	Identifier identifier.Identifier
}

type RetrieveResponse struct {
	Company company.Company
}

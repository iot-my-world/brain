package recordHandler

import (
	"gitlab.com/iotTracker/brain/entity"
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
	Records interface{}
	Total   int
}

type CreateRequest struct {
	Entity entity.Entity
}

type CreateResponse struct {
	Entity entity.Entity
}

type DeleteRequest struct {
	Claims     claims.Claims
	Identifier identifier.Identifier
}

type DeleteResponse struct{}

type UpdateRequest struct {
	Claims     claims.Claims
	Identifier identifier.Identifier
	Entity     entity.Entity
}

type UpdateResponse struct{}

type RetrieveRequest struct {
	Claims     claims.Claims
	Identifier identifier.Identifier
}

type RetrieveResponse struct {
	Entity entity.Entity
}
package genRecordHandler

import (
	"gitlab.com/iotTracker/brain/search/criterion"
	"gitlab.com/iotTracker/brain/search/identifier"
	"gitlab.com/iotTracker/brain/search/query"
	"gitlab.com/iotTracker/brain/security/claims"
)

type RecordHandler interface {
	Create(request *CreateRequest) (*CreateResponse, error)
	Retrieve(request *RetrieveRequest) (*RetrieveResponse, error)
	Update(request *UpdateRequest) (*UpdateResponse, error)
	Delete(request *DeleteRequest) (*DeleteResponse, error)
	Collect(request *CollectRequest) (*CollectResponse, error)
}

type CollectRequest struct {
	Claims   claims.Claims
	Criteria []criterion.Criterion
	Query    query.Query
}

type CollectResponse struct {
	Records []GenEntity
	Total   int
}

type CreateRequest struct {
	Entity GenEntity
}

type CreateResponse struct {
	Entity GenEntity
}

type DeleteRequest struct {
	Claims     claims.Claims
	Identifier identifier.Identifier
}

type DeleteResponse struct {
	Entity GenEntity
}

type UpdateRequest struct {
	Claims     claims.Claims
	Identifier identifier.Identifier
	Entity     GenEntity
}

type UpdateResponse struct {
	Entity GenEntity
}

type RetrieveRequest struct {
	Claims     claims.Claims
	Identifier identifier.Identifier
}

type RetrieveResponse struct {
	Entity GenEntity
}

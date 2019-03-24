package recordHandler

import (
	"gitlab.com/iotTracker/brain/search/criterion"
	"gitlab.com/iotTracker/brain/search/identifier"
	"gitlab.com/iotTracker/brain/search/query"
	"gitlab.com/iotTracker/brain/security/claims"
	"gitlab.com/iotTracker/brain/tracker/reading"
)

type RecordHandler interface {
	Create(request *CreateRequest) (*CreateResponse, error)
	Retrieve(request *RetrieveRequest) (*RetrieveResponse, error)
	Update(request *UpdateRequest) (*UpdateResponse, error)
	Collect(request *CollectRequest) (*CollectResponse, error)
}

type CreateRequest struct {
	Reading reading.Reading
}

type CreateResponse struct {
	Reading reading.Reading
}

type RetrieveRequest struct {
	Claims     claims.Claims
	Identifier identifier.Identifier
}

type RetrieveResponse struct {
	Reading reading.Reading
}

type CollectRequest struct {
	Claims   claims.Claims
	Criteria []criterion.Criterion
	Query    query.Query
}

type CollectResponse struct {
	Records []reading.Reading
	Total   int
}

type UpdateRequest struct {
	Claims     claims.Claims
	Identifier identifier.Identifier
	Reading    reading.Reading
}

type UpdateResponse struct {
	Reading reading.Reading
}

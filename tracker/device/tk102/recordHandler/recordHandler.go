package recordHandler

import (
	"gitlab.com/iotTracker/brain/search/criterion"
	"gitlab.com/iotTracker/brain/search/identifier"
	"gitlab.com/iotTracker/brain/search/query"
	"gitlab.com/iotTracker/brain/security/claims"
	"gitlab.com/iotTracker/brain/tracker/device/tk102"
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
	Records []tk102.TK102
	Total   int
}

type CreateRequest struct {
	TK102 tk102.TK102
}

type CreateResponse struct {
	TK102 tk102.TK102
}

type DeleteRequest struct {
	Identifier identifier.Identifier
}

type DeleteResponse struct {
	TK102 tk102.TK102
}

type UpdateRequest struct {
	Claims     claims.Claims
	Identifier identifier.Identifier
	TK102      tk102.TK102
}

type UpdateResponse struct {
	TK102 tk102.TK102
}

type RetrieveRequest struct {
	Claims     claims.Claims
	Identifier identifier.Identifier
}

type RetrieveResponse struct {
	TK102 tk102.TK102
}

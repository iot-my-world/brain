package recordHandler

import (
	"gitlab.com/iotTracker/brain/party/user"
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

type CreateRequest struct {
	User user.User
}

type CreateResponse struct {
	User user.User
}

type DeleteRequest struct {
	Identifier identifier.Identifier
}

type DeleteResponse struct {
	User user.User
}

type UpdateRequest struct {
	Claims     claims.Claims
	Identifier identifier.Identifier
	User       user.User
}

type UpdateResponse struct {
	User user.User
}

type RetrieveRequest struct {
	Claims     claims.Claims
	Identifier identifier.Identifier
}

type RetrieveResponse struct {
	User user.User
}

type CollectRequest struct {
	Claims   claims.Claims
	Criteria []criterion.Criterion
	Query    query.Query
}

type CollectResponse struct {
	Records []user.User
	Total   int
}

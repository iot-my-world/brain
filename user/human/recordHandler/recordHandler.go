package recordHandler

import (
	"github.com/iot-my-world/brain/search/criterion"
	"github.com/iot-my-world/brain/search/identifier"
	"github.com/iot-my-world/brain/search/query"
	"github.com/iot-my-world/brain/security/claims"
	humanUser "github.com/iot-my-world/brain/user/human"
)

type RecordHandler interface {
	Create(request *CreateRequest) (*CreateResponse, error)
	Retrieve(request *RetrieveRequest) (*RetrieveResponse, error)
	Update(request *UpdateRequest) (*UpdateResponse, error)
	Delete(request *DeleteRequest) (*DeleteResponse, error)
	Collect(request *CollectRequest) (*CollectResponse, error)
}

type CreateRequest struct {
	User humanUser.User
}

type CreateResponse struct {
	User humanUser.User
}

type DeleteRequest struct {
	Identifier identifier.Identifier
}

type DeleteResponse struct {
	User humanUser.User
}

type UpdateRequest struct {
	Claims     claims.Claims
	Identifier identifier.Identifier
	User       humanUser.User
}

type UpdateResponse struct {
	User humanUser.User
}

type RetrieveRequest struct {
	Claims     claims.Claims
	Identifier identifier.Identifier
}

type RetrieveResponse struct {
	User humanUser.User
}

type CollectRequest struct {
	Claims   claims.Claims
	Criteria []criterion.Criterion
	Query    query.Query
}

type CollectResponse struct {
	Records []humanUser.User
	Total   int
}

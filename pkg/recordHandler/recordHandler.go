package recordHandler

import (
	"github.com/iot-my-world/brain/pkg/entity"
	"github.com/iot-my-world/brain/pkg/search/criterion"
	"github.com/iot-my-world/brain/pkg/search/identifier"
	"github.com/iot-my-world/brain/pkg/search/query"
	"github.com/iot-my-world/brain/security/claims"
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

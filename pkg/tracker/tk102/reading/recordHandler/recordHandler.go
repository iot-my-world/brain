package recordHandler

import (
	"github.com/iot-my-world/brain/pkg/search/criterion"
	"github.com/iot-my-world/brain/pkg/search/identifier"
	"github.com/iot-my-world/brain/pkg/search/query"
	"github.com/iot-my-world/brain/pkg/security/claims"
	reading2 "github.com/iot-my-world/brain/pkg/tracker/tk102/reading"
)

type RecordHandler interface {
	Create(request *CreateRequest) (*CreateResponse, error)
	CreateBulk(request *CreateBulkRequest) (*CreateBulkResponse, error)
	Retrieve(request *RetrieveRequest) (*RetrieveResponse, error)
	Update(request *UpdateRequest) (*UpdateResponse, error)
	Collect(request *CollectRequest) (*CollectResponse, error)
}

type CreateRequest struct {
	Reading reading2.Reading
}

type CreateResponse struct {
	Reading reading2.Reading
}

type CreateBulkRequest struct {
	Readings []reading2.Reading
}

type CreateBulkResponse struct {
	Readings []reading2.Reading
}

type RetrieveRequest struct {
	Claims     claims.Claims
	Identifier identifier.Identifier
}

type RetrieveResponse struct {
	Reading reading2.Reading
}

type CollectRequest struct {
	Claims   claims.Claims
	Criteria []criterion.Criterion
	Query    query.Query
}

type CollectResponse struct {
	Records []reading2.Reading
	Total   int
}

type UpdateRequest struct {
	Claims     claims.Claims
	Identifier identifier.Identifier
	Reading    reading2.Reading
}

type UpdateResponse struct {
	Reading reading2.Reading
}

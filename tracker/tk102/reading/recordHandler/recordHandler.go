package recordHandler

import (
	"github.com/iot-my-world/brain/search/criterion"
	"github.com/iot-my-world/brain/search/identifier"
	"github.com/iot-my-world/brain/search/query"
	"github.com/iot-my-world/brain/security/claims"
	"github.com/iot-my-world/brain/tracker/tk102/reading"
)

type RecordHandler interface {
	Create(request *CreateRequest) (*CreateResponse, error)
	CreateBulk(request *CreateBulkRequest) (*CreateBulkResponse, error)
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

type CreateBulkRequest struct {
	Readings []reading.Reading
}

type CreateBulkResponse struct {
	Readings []reading.Reading
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

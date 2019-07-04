package recordHandler

import (
	"github.com/iot-my-world/brain/pkg/search/criterion"
	"github.com/iot-my-world/brain/pkg/search/identifier"
	"github.com/iot-my-world/brain/pkg/search/query"
	"github.com/iot-my-world/brain/pkg/security/claims"
	tk1022 "github.com/iot-my-world/brain/pkg/tracker/tk102"
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
	Records []tk1022.TK102
	Total   int
}

type CreateRequest struct {
	TK102 tk1022.TK102
}

type CreateResponse struct {
	TK102 tk1022.TK102
}

type DeleteRequest struct {
	Identifier identifier.Identifier
}

type DeleteResponse struct {
	TK102 tk1022.TK102
}

type UpdateRequest struct {
	Claims     claims.Claims
	Identifier identifier.Identifier
	TK102      tk1022.TK102
}

type UpdateResponse struct {
	TK102 tk1022.TK102
}

type RetrieveRequest struct {
	Claims     claims.Claims
	Identifier identifier.Identifier
}

type RetrieveResponse struct {
	TK102 tk1022.TK102
}

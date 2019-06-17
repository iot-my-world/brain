package recordHandler

import (
	"github.com/iot-my-world/brain/api"
	"github.com/iot-my-world/brain/party/system"
	"github.com/iot-my-world/brain/search/criterion"
	"github.com/iot-my-world/brain/search/identifier"
	"github.com/iot-my-world/brain/search/query"
	"github.com/iot-my-world/brain/security/claims"
	"github.com/iot-my-world/brain/validate/reasonInvalid"
)

type RecordHandler interface {
	Create(request *CreateRequest) (*CreateResponse, error)
	Retrieve(request *RetrieveRequest) (*RetrieveResponse, error)
	Update(request *UpdateRequest) (*UpdateResponse, error)
	Delete(request *DeleteRequest) (*DeleteResponse, error)
	Validate(request *ValidateRequest) (*ValidateResponse, error)
	Collect(request *CollectRequest) (*CollectResponse, error)
}

const Create api.Method = "Create"
const Retrieve api.Method = "Retrieve"
const Update api.Method = "Update"
const Delete api.Method = "Delete"
const Validate api.Method = "Validate"

type CollectRequest struct {
	Claims   claims.Claims
	Criteria []criterion.Criterion
	Query    query.Query
}

type CollectResponse struct {
	Records []system.System
	Total   int
}

type ValidateRequest struct {
	System system.System
	Method api.Method
}

type ValidateResponse struct {
	ReasonsInvalid []reasonInvalid.ReasonInvalid
}

type CreateRequest struct {
	System system.System
}

type CreateResponse struct {
	System system.System
}

type DeleteRequest struct {
	Claims     claims.Claims
	Identifier identifier.Identifier
}

type DeleteResponse struct {
	System system.System
}

type UpdateRequest struct {
	Claims     claims.Claims
	Identifier identifier.Identifier
	System     system.System
}

type UpdateResponse struct {
	System system.System
}

type RetrieveRequest struct {
	Claims     claims.Claims
	Identifier identifier.Identifier
}

type RetrieveResponse struct {
	System system.System
}

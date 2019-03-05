package recordHandler

import (
	"gitlab.com/iotTracker/brain/api"
	"gitlab.com/iotTracker/brain/party/client"
	"gitlab.com/iotTracker/brain/search/criterion"
	"gitlab.com/iotTracker/brain/search/identifier"
	"gitlab.com/iotTracker/brain/search/query"
	"gitlab.com/iotTracker/brain/security/claims"
	"gitlab.com/iotTracker/brain/validate/reasonInvalid"
)

type RecordHandler interface {
	Create(request *CreateRequest, response *CreateResponse) error
	Retrieve(request *RetrieveRequest, response *RetrieveResponse) error
	Update(request *UpdateRequest, response *UpdateResponse) error
	Delete(request *DeleteRequest, response *DeleteResponse) error
	Validate(request *ValidateRequest, response *ValidateResponse) error
	Collect(request *CollectRequest, response *CollectResponse) error
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
	Records []client.Client
	Total   int
}

type ValidateRequest struct {
	Client client.Client
	Method api.Method
}

type ValidateResponse struct {
	ReasonsInvalid []reasonInvalid.ReasonInvalid
}

type CreateRequest struct {
	Claims claims.Claims
	Client client.Client
}

type CreateResponse struct {
	Client client.Client
}

type DeleteRequest struct {
	Claims     claims.Claims
	Identifier identifier.Identifier
}

type DeleteResponse struct {
	Client client.Client
}

type UpdateRequest struct {
	Claims     claims.Claims
	Identifier identifier.Identifier
	Client     client.Client
}

type UpdateResponse struct {
	Client client.Client
}

type RetrieveRequest struct {
	Claims     claims.Claims
	Identifier identifier.Identifier
}

type RetrieveResponse struct {
	Client client.Client
}

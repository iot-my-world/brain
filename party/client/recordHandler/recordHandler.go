package recordHandler

import (
	"gitlab.com/iotTracker/brain/search/identifier"
	"gitlab.com/iotTracker/brain/validate/reasonInvalid"
	"gitlab.com/iotTracker/brain/party/client"
)

type RecordHandler interface {
	Create(request *CreateRequest, response *CreateResponse) error
	Retrieve(request *RetrieveRequest, response *RetrieveResponse) error
	Update(request *UpdateRequest, response *UpdateResponse) error
	Delete(request *DeleteRequest, response *DeleteResponse) error
	Validate(request *ValidateRequest, response *ValidateResponse) error
}

type ValidateRequest struct {
	Client client.Client
}

type ValidateResponse struct {
	ReasonsInvalid []reasonInvalid.ReasonInvalid
}

type CreateRequest struct {
	Client client.Client
}

type CreateResponse struct {
	Client client.Client
}

type DeleteRequest struct {
	Identifier identifier.Identifier
}

type DeleteResponse struct {
	Client client.Client
}

type UpdateRequest struct {
	Identifier identifier.Identifier
	Client     client.Client
}

type UpdateResponse struct {
	Client client.Client
}

type RetrieveRequest struct {
	Identifier identifier.Identifier
}

type RetrieveResponse struct {
	Client client.Client
}

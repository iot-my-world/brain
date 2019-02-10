package client

import (
	"gitlab.com/iotTracker/brain/search/identifier"
	"gitlab.com/iotTracker/brain/validate/reasonInvalid"
)

type RecordHandler interface {
	Create(request *CreateRequest, response *CreateResponse) error
	Retrieve(request *RetrieveRequest, response *RetrieveResponse) error
	Update(request *UpdateRequest, response *UpdateResponse) error
	Delete(request *DeleteRequest, response *DeleteResponse) error
	Validate(request *ValidateRequest, response *ValidateResponse) error
}

type ValidateRequest struct {
	Client Client
}

type ValidateResponse struct {
	ReasonsInvalid []reasonInvalid.ReasonInvalid
}

type CreateRequest struct {
	Client Client
}

type CreateResponse struct {
	Client Client
}

type DeleteRequest struct {
	Identifier identifier.Identifier
}

type DeleteResponse struct {
	Client Client
}

type UpdateRequest struct {
	Identifier identifier.Identifier
	Client     Client
}

type UpdateResponse struct {
	Client Client
}

type RetrieveRequest struct {
	Identifier identifier.Identifier
}

type RetrieveResponse struct {
	Client Client
}

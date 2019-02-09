package client

import (
	"gitlab.com/iotTracker/brain/search"
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
	Identifier search.Identifier
}

type DeleteResponse struct {
	Client Client
}

type UpdateRequest struct {
	Identifier search.Identifier
	Client     Client
}

type UpdateResponse struct {
	Client Client
}

type RetrieveRequest struct {
	Identifier search.Identifier
}

type RetrieveResponse struct {
	Client Client
}

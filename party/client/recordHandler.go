package client

import (
	"gitlab.com/iotTracker/brain/search"
	"gitlab.com/iotTracker/brain/validate"
)

type RecordHandler interface {
	Create(request *CreateRequest, response *CreateResponse) error
	Retrieve(request *RetrieveRequest, response *RetrieveResponse) error
	Update(request *UpdateRequest, response *UpdateResponse) error
	Delete(request *DeleteRequest, response *DeleteResponse) error
	Validate(request *ValidateRequest, response *ValidateResponse) error
}

type ValidateRequest struct {
	Client Client `json:"client"`
}

type ValidateResponse struct {
	ReasonsInvalid []validate.ReasonInvalid `json:"ReasonsInvalid"`
}

type CreateRequest struct {
	Client Client `json:"client"`
}

type CreateResponse struct {
	Client Client `json:"client"`
}

type DeleteRequest struct {
	Identifier search.Identifier `json:"identifier"`
}

type DeleteResponse struct {
	Client Client `json:"client"`
}

type UpdateRequest struct {
	Identifier search.Identifier `json:"identifier"`
	Client     Client            `json:"client"`
}

type UpdateResponse struct {
	Client Client `json:"client"`
}

type RetrieveRequest struct {
	Identifier search.Identifier `json:"identifier"`
}

type RetrieveResponse struct {
	Client Client `json:"client" bson:"client"`
}

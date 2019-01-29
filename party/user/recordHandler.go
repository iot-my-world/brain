package user

import (
	"gitlab.com/iotTracker/brain/search"
	"gitlab.com/iotTracker/brain/party"
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
	User party.User `json:"user"`
}

type ValidateResponse struct {
	ReasonsInvalid []validate.ReasonInvalid `json:"ReasonsInvalid"`
}

type CreateRequest struct {
	NewUser party.NewUser `json:"newUser"`
}

type CreateResponse struct {
	User party.User `json:"user"`
}

type DeleteRequest struct {
	Username string `json:"username" bson:"username"`
}

type DeleteResponse struct {
	User party.User `json:"user"`
}

type UpdateRequest struct {
	Identifier search.Identifier `json:"identifier"`
	User       party.User        `json:"user"`
}

type UpdateResponse struct {
	User party.User `json:"user"`
}

type RetrieveRequest struct {
	Identifier search.Identifier
}

type RetrieveResponse struct {
	User party.User `json:"user" bson:"user"`
}

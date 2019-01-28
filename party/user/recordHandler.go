package user

import (
	"gitlab.com/iotTracker/brain/search"
	"gitlab.com/iotTracker/brain/party"
)

type RecordHandler interface {
	Create(request *CreateRequest, response *CreateResponse) error
	Retrieve(request *RetrieveRequest, response *RetrieveResponse) error
	RetrieveAll(request *RetrieveAllRequest, response *RetrieveAllResponse) error
	Update(request *UpdateRequest, response *UpdateResponse) error
	Delete(request *DeleteRequest, response *DeleteResponse) error
}

type CreateRequest struct {
	NewUser party.NewUser `json:"newUser"`
}

type CreateResponse struct {
	User party.User `json:"user"`
}

type RetrieveAllRequest struct {
}

type RetrieveAllResponse struct {
	UserRecords []party.User `json:"userRecords" bson:"userRecords"`
}

type DeleteRequest struct {
	Username string `json:"username" bson:"username"`
}

type DeleteResponse struct {
}

type UpdateRequest struct {
	UpdatedUser party.User `json:"updatedUser"`
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

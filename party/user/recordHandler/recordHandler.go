package recordHandler

import (
	"gitlab.com/iotTracker/brain/search"
	"gitlab.com/iotTracker/brain/validate/reasonInvalid"
	"gitlab.com/iotTracker/brain/party/user"
	"gitlab.com/iotTracker/brain/api"
)

type RecordHandler interface {
	Create(request *CreateRequest, response *CreateResponse) error
	Retrieve(request *RetrieveRequest, response *RetrieveResponse) error
	Update(request *UpdateRequest, response *UpdateResponse) error
	Delete(request *DeleteRequest, response *DeleteResponse) error
	Validate(request *ValidateRequest, response *ValidateResponse) error
	ChangePassword(request *ChangePasswordRequest, response *ChangePasswordResponse) error
}

const Create api.Method = "Create"
const Retrieve api.Method = "Retrieve"
const Update api.Method = "Update"
const Delete api.Method = "Delete"
const Validate api.Method = "Validate"

type ValidateRequest struct {
	User user.User
}

type ValidateResponse struct {
	ReasonsInvalid []reasonInvalid.ReasonInvalid
}

type CreateRequest struct {
	User user.User
}

type CreateResponse struct {
	User user.User
}

type DeleteRequest struct {
	Identifier search.Identifier
}

type DeleteResponse struct {
	User user.User
}

type UpdateRequest struct {
	Identifier search.Identifier
	User       user.User
}

type UpdateResponse struct {
	User user.User
}

type RetrieveRequest struct {
	Identifier search.Identifier
}

type RetrieveResponse struct {
	User user.User
}

type ChangePasswordRequest struct {
	Identifier  search.Identifier
	NewPassword string
}

type ChangePasswordResponse struct {
	User user.User
}

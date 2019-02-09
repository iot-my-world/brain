package user

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
	ChangePassword(request *ChangePasswordRequest, response *ChangePasswordResponse) error
}

type ValidateRequest struct {
	User User
}

type ValidateResponse struct {
	ReasonsInvalid []reasonInvalid.ReasonInvalid
}

type CreateRequest struct {
	User User
}

type CreateResponse struct {
	User User
}

type DeleteRequest struct {
	Identifier search.Identifier
}

type DeleteResponse struct {
	User User
}

type UpdateRequest struct {
	Identifier search.Identifier
	User       User
}

type UpdateResponse struct {
	User User
}

type RetrieveRequest struct {
	Identifier search.Identifier
}

type RetrieveResponse struct {
	User User
}

type ChangePasswordRequest struct {
	Identifier  search.Identifier
	NewPassword string
}

type ChangePasswordResponse struct {
	User User
}

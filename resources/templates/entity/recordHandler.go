package company

import (
	"gitlab.com/iotTracker/brain/search"
	"gitlab.com/iotTracker/brain/validate"
	"gitlab.com/iotTracker/brain/security/claims"
)

type RecordHandler interface {
	Create(request *CreateRequest, response *CreateResponse) error
	Retrieve(request *RetrieveRequest, response *RetrieveResponse) error
	Update(request *UpdateRequest, response *UpdateResponse) error
	Delete(request *DeleteRequest, response *DeleteResponse) error
	Validate(request *ValidateRequest, response *ValidateResponse) error
}

type ValidateRequest struct {
	Company Company
}

type ValidateResponse struct {
	ReasonsInvalid []validate.ReasonInvalid
}

type CreateRequest struct {
	Company Company
}

type CreateResponse struct {
	Company Company
}

type DeleteRequest struct {
	Identifier identifier.Identifier
}

type DeleteResponse struct {
	Company Company
}

type UpdateRequest struct {
	Claims     claims.Claims
	Identifier identifier.Identifier
	Company    Company
}

type UpdateResponse struct {
	Company Company
}

type RetrieveRequest struct {
	Claims     claims.Claims
	Identifier identifier.Identifier
}

type RetrieveResponse struct {
	Company Company
}

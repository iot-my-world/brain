package recordHandler

import (
	"gitlab.com/iotTracker/brain/search/identifier"
	"gitlab.com/iotTracker/brain/security/role"
)

type RecordHandler interface {
	Create(request *CreateRequest, response *CreateResponse) error
	Retrieve(request *RetrieveRequest, response *RetrieveResponse) error
	Update(request *UpdateRequest, response *UpdateResponse) error
}

type CreateRequest struct {
	Role role.Role `json:"role"`
}

type CreateResponse struct {
}

type RetrieveRequest struct {
	Identifier identifier.Identifier `json:"identifier"`
}

type RetrieveResponse struct {
	Role role.Role `json:"role"`
}

type UpdateRequest struct {
	Identifier identifier.Identifier `json:"identifier"`
	Role       role.Role             `json:"role"`
}

type UpdateResponse struct {
}
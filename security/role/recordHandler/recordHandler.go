package recordHandler

import (
	"github.com/iot-my-world/brain/search/identifier"
	"github.com/iot-my-world/brain/security/role"
)

type RecordHandler interface {
	Create(request *CreateRequest) (*CreateResponse, error)
	Retrieve(request *RetrieveRequest) (*RetrieveResponse, error)
	Update(request *UpdateRequest) (*UpdateResponse, error)
}

type CreateRequest struct {
	Role role.Role
}

type CreateResponse struct {
	Role role.Role
}

type RetrieveRequest struct {
	Identifier identifier.Identifier
}

type RetrieveResponse struct {
	Role role.Role
}

type UpdateRequest struct {
	Identifier identifier.Identifier
	Role       role.Role
}

type UpdateResponse struct {
	Role role.Role
}

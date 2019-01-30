package role

import (
	"gitlab.com/iotTracker/brain/security"
	"gitlab.com/iotTracker/brain/search"
)

type RecordHandler interface {
	Create(request *CreateRequest, response *CreateResponse) error
	Retrieve(request *RetrieveRequest, response *RetrieveResponse) error
	Update(request *UpdateRequest, response *UpdateResponse) error
}

type CreateRequest struct {
	Role security.Role `json:"role"`
}

type CreateResponse struct {
}

type RetrieveRequest struct {
	Identifier search.Identifier `json:"identifier"`
}

type RetrieveResponse struct {
	Role security.Role `json:"role"`
}

type UpdateRequest struct {
	Role security.Role `json:"role"`
}

type UpdateResponse struct {
}

package systemRole

import "gitlab.com/iotTracker/brain/security"

type RecordHandler interface {
	Create(request *CreateRequest, response *CreateResponse) error
	Retrieve(request *RetrieveRequest, response *RetrieveResponse) error
	Update(request *UpdateRequest, response *UpdateResponse) error
}

type CreateRequest struct {
	SystemRole security.SystemRole `json:"role" bson:"role"`
}

type CreateResponse struct {
}

type RetrieveRequest struct {
	Name string `json:"name" bson:"name"`
}

type RetrieveResponse struct {
	SystemRole security.SystemRole `json:"systemRole" bson:"systemRole"`
}

type UpdateRequest struct {
	SystemRole security.SystemRole `json:"systemRole" bson:"systemRole"`
}

type UpdateResponse struct {
}

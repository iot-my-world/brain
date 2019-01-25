package systemRole

type RecordHandler interface {
	Create(request *CreateRequest, response *CreateResponse) error
	Retrieve(request *RetrieveRequest, response *RetrieveResponse) error
	Update(request *UpdateRequest, response *UpdateResponse) error
}

type CreateRequest struct {
	SystemRole SystemRole `json:"role" bson:"role"`
}

type CreateResponse struct {
}

type RetrieveRequest struct {
	Name string `json:"name" bson:"name"`
}

type RetrieveResponse struct {
	SystemRole SystemRole `json:"systemRole" bson:"systemRole"`
}

type UpdateRequest struct {
	SystemRole SystemRole `json:"systemRole" bson:"systemRole"`
}

type UpdateResponse struct {
}

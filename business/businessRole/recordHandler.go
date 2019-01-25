package businessRole

type RecordHandler interface {
	Create(request *CreateRequest, response *CreateResponse) error
	RetrieveAll(request *RetrieveAllRequest, response *RetrieveAllResponse) error
	Update(request *UpdateRequest, response *UpdateResponse) error
	Delete(request *DeleteRequest, response *DeleteResponse) error
}

type CreateRequest struct {
	BusinessRole  `json:"businessRole" bson:"businessRole"`
}

type CreateResponse struct {
	BusinessRole BusinessRole `json:"businessRole" bson:"businessRole"`
}

type RetrieveAllRequest struct {}

type RetrieveAllResponse struct {
	Records []BusinessRole `json:"records"`
}

type UpdateRequest struct {
	BusinessRole BusinessRole `json:"businessRole" bson:"businessRole"`
}

type UpdateResponse struct {
	BusinessRole BusinessRole `json:"businessRole" bson:"businessRole"`
}

type DeleteRequest struct {
	BusinessRole BusinessRole `json:"businessRole" bson:"businessRole"`
}

type DeleteResponse struct {
	BusinessRole BusinessRole `json:"businessRole" bson:"businessRole"`
}
package ship

type RecordHandler interface {
	Create(request *CreateRequest, response *CreateResponse) error
	RetrieveAll(request *RetrieveAllRequest, response *RetrieveAllResponse) error
	Update(request *UpdateRequest, response *UpdateResponse) error
	Delete(request *DeleteRequest, response *DeleteResponse) error
}

type CreateRequest struct {
	Ship  `json:"ship" bson:"ship"`
}

type CreateResponse struct {
	Ship Ship `json:"ship" bson:"ship"`
}

type RetrieveAllRequest struct {}

type RetrieveAllResponse struct {
	Records []Ship `json:"records"`
}

type UpdateRequest struct {
	Ship Ship `json:"ship" bson:"ship"`
}

type UpdateResponse struct {
	Ship Ship `json:"ship" bson:"ship"`
}

type DeleteRequest struct {
	Ship Ship `json:"ship" bson:"ship"`
}

type DeleteResponse struct {
	Ship Ship `json:"ship" bson:"ship"`
}
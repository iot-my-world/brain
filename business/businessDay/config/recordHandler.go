package config

type RecordHandler interface {
	Create(request *CreateRequest, response *CreateResponse) error
	Update(request *UpdateRequest, response *UpdateResponse) error
	Retrieve(RetrieveRequest *RetrieveRequest, response *RetrieveResponse) error
}

type CreateRequest struct {
	BusinessDayConfig Config `json:"businessDayConfig"`
}

type CreateResponse struct {
	BusinessDayConfig Config `json:"businessDayConfig"`
}

type UpdateRequest struct {
	BusinessDayConfig Config `json:"businessDayConfig"`
}

type UpdateResponse struct {
	BusinessDayConfig Config `json:"businessDayConfig"`
}

type RetrieveRequest struct {}

type RetrieveResponse struct {
	BusinessDayConfig Config `json:"businessDayConfig"`
}


package systemRole

import (
	"net/http"
)

type service struct{
	RecordHandler
}

func NewService(recordHandler RecordHandler) *service {
	return &service{
		recordHandler,
	}
}

func (s *service) Create(r *http.Request, request *CreateRequest, response *CreateResponse) error {
	return s.RecordHandler.Create(request, response)
}

func (s *service) Retrieve(r *http.Request, request *RetrieveRequest, response *RetrieveResponse) error {
	return s.RecordHandler.Retrieve(request, response)
}

func (s *service) Update(r *http.Request, request *UpdateRequest, response *UpdateResponse) error {
	return s.RecordHandler.Update(request, response)
}

package ship

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

func (s *service) RetrieveAll(r *http.Request, request *RetrieveAllRequest, response *RetrieveAllResponse) error {
	return s.RecordHandler.RetrieveAll(request, response)
}

func (s *service) Update(r *http.Request, request *UpdateRequest, response *UpdateResponse) error {
	return s.RecordHandler.Update(request, response)
}

func (s *service) Delete(r *http.Request, request *DeleteRequest, response *DeleteResponse) error {
	return s.RecordHandler.Delete(request, response)
}

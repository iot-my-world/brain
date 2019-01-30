package user

import (
	"net/http"
)

type serviceAdaptor struct {
	RecordHandler
}

func NewServiceAdaptor(recordHandler RecordHandler) *serviceAdaptor {
	return &serviceAdaptor{
		recordHandler,
	}
}

func (s *serviceAdaptor) Create(r *http.Request, request *CreateRequest, response *CreateResponse) error {
	return s.RecordHandler.Create(request, response)
}

func (s *serviceAdaptor) Update(r *http.Request, request *UpdateRequest, response *UpdateResponse) error {
	return s.RecordHandler.Update(request, response)
}

func (s *serviceAdaptor) Delete(r *http.Request, request *DeleteRequest, response *DeleteResponse) error {
	return s.RecordHandler.Delete(request, response)
}

package role

import (
	"gitlab.com/iotTracker/brain/security/role"
	"net/http"
)

type adaptor struct {
	role.RecordHandler
}

func New(recordHandler role.RecordHandler) *adaptor {
	return &adaptor{
		RecordHandler: recordHandler,
	}
}

func (s *adaptor) Create(r *http.Request, request *role.CreateRequest, response *role.CreateResponse) error {
	return s.RecordHandler.Create(request, response)
}

func (s *adaptor) Retrieve(r *http.Request, request *role.RetrieveRequest, response *role.RetrieveResponse) error {
	return s.RecordHandler.Retrieve(request, response)
}

func (s *adaptor) Update(r *http.Request, request *role.UpdateRequest, response *role.UpdateResponse) error {
	return s.RecordHandler.Update(request, response)
}

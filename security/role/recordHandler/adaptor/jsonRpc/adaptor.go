package role

import (
	"net/http"
	"gitlab.com/iotTracker/brain/security/role"
)

type service struct{
	role.RecordHandler
}

func New(recordHandler role.RecordHandler) *service {
	return &service{
		RecordHandler: recordHandler,
	}
}

func (s *service) Create(r *http.Request, request *role.CreateRequest, response *role.CreateResponse) error {
	return s.RecordHandler.Create(request, response)
}

func (s *service) Retrieve(r *http.Request, request *role.RetrieveRequest, response *role.RetrieveResponse) error {
	return s.RecordHandler.Retrieve(request, response)
}

func (s *service) Update(r *http.Request, request *role.UpdateRequest, response *role.UpdateResponse) error {
	return s.RecordHandler.Update(request, response)
}

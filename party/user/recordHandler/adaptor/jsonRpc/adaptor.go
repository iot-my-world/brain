package user

import (
	"net/http"
	"gitlab.com/iotTracker/brain/search/identifier/id"
	"gitlab.com/iotTracker/brain/party/user"
)

type adaptor struct {
	RecordHandler user.RecordHandler
}

func New(recordHandler user.RecordHandler) *adaptor {
	return &adaptor{
		recordHandler,
	}
}

func (s *adaptor) Create(r *http.Request, request *user.CreateRequest, response *user.CreateResponse) error {
	return s.RecordHandler.Create(request, response)
}

func (s *adaptor) Retrieve(r *http.Request, request *user.RetrieveRequest, response *user.RetrieveResponse) error {
	return s.RecordHandler.Retrieve(&user.RetrieveRequest{Identifier: id.Identifier{Id: "1234"}}, response)
}

func (s *adaptor) Update(r *http.Request, request *user.UpdateRequest, response *user.UpdateResponse) error {
	return s.RecordHandler.Update(request, response)
}

func (s *adaptor) Delete(r *http.Request, request *user.DeleteRequest, response *user.DeleteResponse) error {
	return s.RecordHandler.Delete(request, response)
}

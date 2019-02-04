package user

import (
	"net/http"
	"gitlab.com/iotTracker/brain/search/identifier/id"
	"gitlab.com/iotTracker/brain/party/user"
)

type serviceAdaptor struct {
	RecordHandler user.RecordHandler
}

func New(recordHandler user.RecordHandler) *serviceAdaptor {
	return &serviceAdaptor{
		recordHandler,
	}
}

func (s *serviceAdaptor) Create(r *http.Request, request *user.CreateRequest, response *user.CreateResponse) error {
	return s.RecordHandler.Create(request, response)
}

func (s *serviceAdaptor) Retrieve(r *http.Request, request *user.RetrieveRequest, response *user.RetrieveResponse) error {
	return s.RecordHandler.Retrieve(&user.RetrieveRequest{Identifier: id.Identifier("5c520cdf5402664b2017c0fc")}, response)
}

func (s *serviceAdaptor) Update(r *http.Request, request *user.UpdateRequest, response *user.UpdateResponse) error {
	return s.RecordHandler.Update(request, response)
}

func (s *serviceAdaptor) Delete(r *http.Request, request *user.DeleteRequest, response *user.DeleteResponse) error {
	return s.RecordHandler.Delete(request, response)
}

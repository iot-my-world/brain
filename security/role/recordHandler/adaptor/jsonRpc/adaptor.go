package role

import (
	roleRecordHandler "gitlab.com/iotTracker/brain/security/role/recordHandler"
	"net/http"
)

type adaptor struct {
	roleRecordHandler.RecordHandler
}

func New(recordHandler roleRecordHandler.RecordHandler) *adaptor {
	return &adaptor{
		RecordHandler: recordHandler,
	}
}

func (s *adaptor) Create(r *http.Request, request *roleRecordHandler.CreateRequest, response *roleRecordHandler.CreateResponse) error {
	return s.RecordHandler.Create(request, response)
}

func (s *adaptor) Retrieve(r *http.Request, request *roleRecordHandler.RetrieveRequest, response *roleRecordHandler.RetrieveResponse) error {
	return s.RecordHandler.Retrieve(request, response)
}

func (s *adaptor) Update(r *http.Request, request *roleRecordHandler.UpdateRequest, response *roleRecordHandler.UpdateResponse) error {
	return s.RecordHandler.Update(request, response)
}

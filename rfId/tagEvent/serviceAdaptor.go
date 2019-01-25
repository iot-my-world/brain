package tagEvent

import (
	"net/http"
	"time"
)

type serviceAdaptor struct{
	RecordHandler
}

func NewServiceAdaptor(handler RecordHandler) *serviceAdaptor {
	return &serviceAdaptor{
		RecordHandler: handler,
	}
}

func (s *serviceAdaptor) RFIDTagEvent(r *http.Request, request *RFIDTagEventRequest, response *RFIDTagEventResponse) error {
	request.TagEvent.ReceivedTime = time.Now().UTC().Unix()

	return s.RecordHandler.RFIDTagEvent(request, response)
}

func (s *serviceAdaptor) Retrieve(r *http.Request, request *RetrieveRequest, response *RetrieveResponse) error {
	return nil
}

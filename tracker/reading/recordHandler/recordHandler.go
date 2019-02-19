package recordHandler

import "gitlab.com/iotTracker/brain/tracker/reading"

type RecordHandler interface {
	Create (request *CreateRequest, response *CreateResponse) error
}

type CreateRequest struct {
	Reading reading.Reading
}

type CreateResponse struct {
	Reading reading.Reading
}
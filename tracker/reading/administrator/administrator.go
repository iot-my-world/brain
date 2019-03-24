package administrator

import "gitlab.com/iotTracker/brain/tracker/reading"

type Administrator interface {
	Create(request *CreateRequest) (*CreateResponse, error)
}

type CreateRequest struct {
	Reading reading.Reading
}

type CreateResponse struct {
	Reading reading.Reading
}

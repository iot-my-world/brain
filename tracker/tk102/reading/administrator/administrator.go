package administrator

import "gitlab.com/iotTracker/brain/tracker/tk102/reading"

type Administrator interface {
	Create(request *CreateRequest) (*CreateResponse, error)
	CreateBulk(request *CreateBulkRequest) (*CreateBulkResponse, error)
}

type CreateRequest struct {
	Reading reading.Reading
}

type CreateResponse struct {
	Reading reading.Reading
}

type CreateBulkRequest struct {
	Readings []reading.Reading
}

type CreateBulkResponse struct {
	Readings []reading.Reading
}

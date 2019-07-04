package administrator

import (
	reading2 "github.com/iot-my-world/brain/pkg/tracker/tk102/reading"
)

type Administrator interface {
	Create(request *CreateRequest) (*CreateResponse, error)
	CreateBulk(request *CreateBulkRequest) (*CreateBulkResponse, error)
}

type CreateRequest struct {
	Reading reading2.Reading
}

type CreateResponse struct {
	Reading reading2.Reading
}

type CreateBulkRequest struct {
	Readings []reading2.Reading
}

type CreateBulkResponse struct {
	Readings []reading2.Reading
}

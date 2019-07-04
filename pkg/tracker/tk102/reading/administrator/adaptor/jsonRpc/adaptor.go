package jsonRpc

import (
	reading2 "github.com/iot-my-world/brain/pkg/tracker/tk102/reading"
	"github.com/iot-my-world/brain/pkg/tracker/tk102/reading/administrator"
	"net/http"
)

type Adaptor struct {
	administrator administrator.Administrator
}

func New(administrator administrator.Administrator) *Adaptor {
	return &Adaptor{
		administrator: administrator,
	}
}

type CreateRequest struct {
	Reading reading2.Reading `json:"reading"`
}

type CreateResponse struct {
	Reading reading2.Reading `json:"reading"`
}

func (a *Adaptor) Create(r *http.Request, request *CreateRequest, response *CreateResponse) error {

	createResponse, err := a.administrator.Create(&administrator.CreateRequest{
		Reading: request.Reading,
	})
	if err != nil {
		return err
	}

	response.Reading = createResponse.Reading

	return nil
}

type CreateBulkRequest struct {
	Readings []reading2.Reading `json:"readings"`
}

type CreateBulkResponse struct {
	Readings []reading2.Reading `json:"readings"`
}

func (a *Adaptor) CreateBulk(r *http.Request, request *CreateBulkRequest, response *CreateBulkResponse) error {

	createResponse, err := a.administrator.CreateBulk(&administrator.CreateBulkRequest{
		Readings: request.Readings,
	})
	if err != nil {
		return err
	}

	response.Readings = createResponse.Readings

	return nil
}

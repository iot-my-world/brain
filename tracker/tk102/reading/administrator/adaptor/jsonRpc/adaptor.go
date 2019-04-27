package jsonRpc

import (
	"gitlab.com/iotTracker/brain/tracker/tk102/reading"
	readingAdministrator "gitlab.com/iotTracker/brain/tracker/tk102/reading/administrator"
	"net/http"
)

type Adaptor struct {
	administrator readingAdministrator.Administrator
}

func New(administrator readingAdministrator.Administrator) *Adaptor {
	return &Adaptor{
		administrator: administrator,
	}
}

type CreateRequest struct {
	Reading reading.Reading `json:"reading"`
}

type CreateResponse struct {
	Reading reading.Reading `json:"reading"`
}

func (a *Adaptor) Create(r *http.Request, request *CreateRequest, response *CreateResponse) error {

	createResponse, err := a.administrator.Create(&readingAdministrator.CreateRequest{
		Reading: request.Reading,
	})
	if err != nil {
		return err
	}

	response.Reading = createResponse.Reading

	return nil
}

type CreateBulkRequest struct {
	Readings []reading.Reading `json:"readings"`
}

type CreateBulkResponse struct {
	Readings []reading.Reading `json:"readings"`
}

func (a *Adaptor) CreateBulk(r *http.Request, request *CreateBulkRequest, response *CreateBulkResponse) error {

	createResponse, err := a.administrator.CreateBulk(&readingAdministrator.CreateBulkRequest{
		Readings: request.Readings,
	})
	if err != nil {
		return err
	}

	response.Readings = createResponse.Readings

	return nil
}

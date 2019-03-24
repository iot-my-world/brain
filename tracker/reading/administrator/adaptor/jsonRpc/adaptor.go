package jsonRpc

import (
	"gitlab.com/iotTracker/brain/tracker/reading"
	readingAdministrator "gitlab.com/iotTracker/brain/tracker/reading/administrator"
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

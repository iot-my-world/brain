package jsonRpc

import (
	"gitlab.com/iotTracker/brain/log"
	wrappedClaims "gitlab.com/iotTracker/brain/security/claims/wrapped"
	zx303Task "gitlab.com/iotTracker/brain/tracker/zx303/task"
	zx303TaskAdministrator "gitlab.com/iotTracker/brain/tracker/zx303/task/administrator"
	"net/http"
)

type Adaptor struct {
	administrator zx303TaskAdministrator.Administrator
}

func New(administrator zx303TaskAdministrator.Administrator) *Adaptor {
	return &Adaptor{
		administrator: administrator,
	}
}

type CreateRequest struct {
	ZX303Task zx303Task.Task `json:"zx303Task"`
}

type CreateResponse struct {
	ZX303Task zx303Task.Task `json:"zx303Task"`
}

func (a *Adaptor) Create(r *http.Request, request *CreateRequest, response *CreateResponse) error {
	claims, err := wrappedClaims.UnwrapClaimsFromContext(r)
	if err != nil {
		log.Warn(err.Error())
		return err
	}

	createResponse, err := a.administrator.Create(&zx303TaskAdministrator.CreateRequest{
		Claims:    claims,
		ZX303Task: request.ZX303Task,
	})
	if err != nil {
		return err
	}

	response.ZX303Task = createResponse.ZX303Task

	return nil
}

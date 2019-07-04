package jsonRpc

import (
	"github.com/iot-my-world/brain/internal/log"
	wrappedClaims "github.com/iot-my-world/brain/pkg/security/claims/wrapped"
	"github.com/iot-my-world/brain/pkg/tracker/zx303/reading/status"
	"github.com/iot-my-world/brain/pkg/tracker/zx303/reading/status/administrator"
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
	ZX303StatusReading status.Reading `json:"zx303StatusReading"`
}

type CreateResponse struct {
	ZX303StatusReading status.Reading `json:"zx303StatusReading"`
}

func (a *Adaptor) Create(r *http.Request, request *CreateRequest, response *CreateResponse) error {
	claims, err := wrappedClaims.UnwrapClaimsFromContext(r)
	if err != nil {
		log.Warn(err.Error())
		return err
	}

	createResponse, err := a.administrator.Create(&administrator.CreateRequest{
		Claims:             claims,
		ZX303StatusReading: request.ZX303StatusReading,
	})
	if err != nil {
		return err
	}

	response.ZX303StatusReading = createResponse.ZX303StatusReading

	return nil
}

package jsonRpc

import (
	"github.com/iot-my-world/brain/log"
	wrappedClaims "github.com/iot-my-world/brain/security/claims/wrapped"
	zx303StatusReading "github.com/iot-my-world/brain/tracker/zx303/reading/status"
	zx303StatusReadingAdministrator "github.com/iot-my-world/brain/tracker/zx303/reading/status/administrator"
	"net/http"
)

type Adaptor struct {
	administrator zx303StatusReadingAdministrator.Administrator
}

func New(administrator zx303StatusReadingAdministrator.Administrator) *Adaptor {
	return &Adaptor{
		administrator: administrator,
	}
}

type CreateRequest struct {
	ZX303StatusReading zx303StatusReading.Reading `json:"zx303StatusReading"`
}

type CreateResponse struct {
	ZX303StatusReading zx303StatusReading.Reading `json:"zx303StatusReading"`
}

func (a *Adaptor) Create(r *http.Request, request *CreateRequest, response *CreateResponse) error {
	claims, err := wrappedClaims.UnwrapClaimsFromContext(r)
	if err != nil {
		log.Warn(err.Error())
		return err
	}

	createResponse, err := a.administrator.Create(&zx303StatusReadingAdministrator.CreateRequest{
		Claims:             claims,
		ZX303StatusReading: request.ZX303StatusReading,
	})
	if err != nil {
		return err
	}

	response.ZX303StatusReading = createResponse.ZX303StatusReading

	return nil
}

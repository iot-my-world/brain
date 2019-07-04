package jsonRpc

import (
	"github.com/iot-my-world/brain/internal/log"
	wrappedClaims "github.com/iot-my-world/brain/pkg/security/claims/wrapped"
	"github.com/iot-my-world/brain/pkg/tracker/zx303/reading/gps"
	"github.com/iot-my-world/brain/pkg/tracker/zx303/reading/gps/administrator"
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
	ZX303GPSReading gps.Reading `json:"zx303GPSReading"`
}

type CreateResponse struct {
	ZX303GPSReading gps.Reading `json:"zx303GPSReading"`
}

func (a *Adaptor) Create(r *http.Request, request *CreateRequest, response *CreateResponse) error {
	claims, err := wrappedClaims.UnwrapClaimsFromContext(r)
	if err != nil {
		log.Warn(err.Error())
		return err
	}

	createResponse, err := a.administrator.Create(&administrator.CreateRequest{
		Claims:          claims,
		ZX303GPSReading: request.ZX303GPSReading,
	})
	if err != nil {
		return err
	}

	response.ZX303GPSReading = createResponse.ZX303GPSReading

	return nil
}

package jsonRpc

import (
	"github.com/iot-my-world/brain/log"
	wrappedClaims "github.com/iot-my-world/brain/security/claims/wrapped"
	zx303GPSReading "github.com/iot-my-world/brain/tracker/zx303/reading/gps"
	zx303GPSReadingAdministrator "github.com/iot-my-world/brain/tracker/zx303/reading/gps/administrator"
	"net/http"
)

type Adaptor struct {
	administrator zx303GPSReadingAdministrator.Administrator
}

func New(administrator zx303GPSReadingAdministrator.Administrator) *Adaptor {
	return &Adaptor{
		administrator: administrator,
	}
}

type CreateRequest struct {
	ZX303GPSReading zx303GPSReading.Reading `json:"zx303GPSReading"`
}

type CreateResponse struct {
	ZX303GPSReading zx303GPSReading.Reading `json:"zx303GPSReading"`
}

func (a *Adaptor) Create(r *http.Request, request *CreateRequest, response *CreateResponse) error {
	claims, err := wrappedClaims.UnwrapClaimsFromContext(r)
	if err != nil {
		log.Warn(err.Error())
		return err
	}

	createResponse, err := a.administrator.Create(&zx303GPSReadingAdministrator.CreateRequest{
		Claims:          claims,
		ZX303GPSReading: request.ZX303GPSReading,
	})
	if err != nil {
		return err
	}

	response.ZX303GPSReading = createResponse.ZX303GPSReading

	return nil
}

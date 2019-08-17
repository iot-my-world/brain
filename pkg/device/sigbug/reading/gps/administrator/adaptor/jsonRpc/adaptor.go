package jsonRpc

import (
	"github.com/iot-my-world/brain/internal/log"
	jsonRpcServiceProvider "github.com/iot-my-world/brain/pkg/api/jsonRpc/service/provider"
	sigbugGPSReading "github.com/iot-my-world/brain/pkg/device/sigbug/reading/gps"
	sigbugGPSReadingAdministrator "github.com/iot-my-world/brain/pkg/device/sigbug/reading/gps/administrator"
	wrappedClaims "github.com/iot-my-world/brain/pkg/security/claims/wrapped"
	"net/http"
)

type adaptor struct {
	sigbugGPSReadingAdministrator sigbugGPSReadingAdministrator.Administrator
}

func New(administrator sigbugGPSReadingAdministrator.Administrator) *adaptor {
	return &adaptor{
		sigbugGPSReadingAdministrator: administrator,
	}
}

func (a *adaptor) Name() jsonRpcServiceProvider.Name {
	return sigbugGPSReadingAdministrator.ServiceProvider
}

func (a *adaptor) MethodRequiresAuthorization(string) bool {
	return true
}

type CreateRequest struct {
	Reading sigbugGPSReading.Reading `json:"reading"`
}

type CreateResponse struct {
	Reading sigbugGPSReading.Reading `json:"reading"`
}

func (a *adaptor) Create(r *http.Request, request *CreateRequest, response *CreateResponse) error {
	claims, err := wrappedClaims.UnwrapClaimsFromContext(r)
	if err != nil {
		log.Warn(err.Error())
		return err
	}

	createResponse, err := a.sigbugGPSReadingAdministrator.Create(&sigbugGPSReadingAdministrator.CreateRequest{
		Claims:  claims,
		Reading: request.Reading,
	})
	if err != nil {
		return err
	}

	response.Reading = createResponse.Reading

	return nil
}

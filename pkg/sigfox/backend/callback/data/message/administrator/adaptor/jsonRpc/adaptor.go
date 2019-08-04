package jsonRpc

import (
	"github.com/iot-my-world/brain/internal/log"
	jsonRpcServiceProvider "github.com/iot-my-world/brain/pkg/api/jsonRpc/service/provider"
	wrappedClaims "github.com/iot-my-world/brain/pkg/security/claims/wrapped"
	sigfoxBackendDataCallbackMessage "github.com/iot-my-world/brain/pkg/sigfox/backend/callback/data/message"
	"github.com/iot-my-world/brain/pkg/sigfox/backend/callback/data/message/administrator"
	"net/http"
)

type adaptor struct {
	administrator administrator.Administrator
}

func New(administrator administrator.Administrator) *adaptor {
	return &adaptor{
		administrator: administrator,
	}
}

func (a *adaptor) Name() jsonRpcServiceProvider.Name {
	return jsonRpcServiceProvider.Name(administrator.ServiceProvider)
}

func (a *adaptor) MethodRequiresAuthorization(string) bool {
	return true
}

type CreateRequest struct {
	Message sigfoxBackendDataCallbackMessage.Message `json:"message"`
}

type CreateResponse struct {
	Message sigfoxBackendDataCallbackMessage.Message `json:"message"`
}

func (a *adaptor) Create(r *http.Request, request *CreateRequest, response *CreateResponse) error {
	claims, err := wrappedClaims.UnwrapClaimsFromContext(r)
	if err != nil {
		log.Warn(err.Error())
		return err
	}

	createResponse, err := a.administrator.Create(&administrator.CreateRequest{
		Claims:  claims,
		Message: request.Message,
	})
	if err != nil {
		return err
	}

	response.Message = createResponse.Message

	return nil
}

package jsonRpc

import (
	"encoding/hex"
	"github.com/iot-my-world/brain/internal/log"
	jsonRpcServiceProvider "github.com/iot-my-world/brain/pkg/api/jsonRpc/service/provider"
	"github.com/iot-my-world/brain/pkg/search/identifier/id"
	sigfoxBackendCallbackDataMessage "github.com/iot-my-world/brain/pkg/sigfox/backend/callback/data/message"
	sigfoxBackendCallbackServer "github.com/iot-my-world/brain/pkg/sigfox/backend/callback/server"
	"net/http"
)

type adaptor struct {
	Server sigfoxBackendCallbackServer.Server
}

func New(
	Server sigfoxBackendCallbackServer.Server,
) *adaptor {
	return &adaptor{
		Server: Server,
	}
}

func (a *adaptor) Name() jsonRpcServiceProvider.Name {
	return jsonRpcServiceProvider.Name(sigfoxBackendCallbackServer.ServiceProvider)
}

func (a *adaptor) MethodRequiresAuthorization(string) bool {
	return true
}

type HandleDataMessageRequest struct {
	DeviceIdentifier id.Identifier `json:"deviceIdentifier"`
	Data             string        `json:"data"`
}

type HandleDataMessageResponse struct {
}

func (a *adaptor) HandleDataMessage(r *http.Request, request *HandleDataMessageRequest, response *HandleDataMessageResponse) error {
	messageData, err := hex.DecodeString(request.Data)
	if err != nil {
		log.Error(err.Error())
		return err
	}

	if _, err := a.Server.HandleDataMessage(&sigfoxBackendCallbackServer.HandleDataMessageRequest{
		Message: sigfoxBackendCallbackDataMessage.Message{
			DeviceIdentifier: request.DeviceIdentifier,
			Data:             messageData,
		},
	}); err != nil {
		return err
	}

	return nil
}

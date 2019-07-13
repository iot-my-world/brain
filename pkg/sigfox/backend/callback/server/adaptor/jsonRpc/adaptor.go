package jsonRpc

import (
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

type HandleDataMessageRequest struct {
	Device string `json:"device"`
	Data   []byte `json:"data"`
}

type HandleDataMessageResponse struct {
}

func (a *adaptor) HandleDataMessage(r *http.Request, request *HandleDataMessageRequest, response *HandleDataMessageResponse) error {
	if _, err := a.Server.HandleDataMessage(&sigfoxBackendCallbackServer.HandleDataMessageRequest{
		Message: sigfoxBackendCallbackDataMessage.Message{
			DeviceId: request.Device,
			Data:     request.Data,
		},
	}); err != nil {
		return err
	}

	return nil
}

package basic

import (
	"github.com/iot-my-world/brain/internal/log"
	sigfoxBackendDataCallbackMessageAdministrator "github.com/iot-my-world/brain/pkg/sigfox/backend/callback/data/message/administrator"
	sigfoxBackendDataMessageHandler "github.com/iot-my-world/brain/pkg/sigfox/backend/callback/data/message/handler"
	sigfoxBackendCallbackServer "github.com/iot-my-world/brain/pkg/sigfox/backend/callback/server"
	sigfoxBackendCallbackServerException "github.com/iot-my-world/brain/pkg/sigfox/backend/callback/server/exception"
)

type server struct {
	handlers                                      []sigfoxBackendDataMessageHandler.Handler
	sigfoxBackendDataCallbackMessageAdministrator sigfoxBackendDataCallbackMessageAdministrator.Administrator
}

func New(
	sigfoxBackendDataCallbackMessageAdministrator sigfoxBackendDataCallbackMessageAdministrator.Administrator,
	handlers []sigfoxBackendDataMessageHandler.Handler,
) sigfoxBackendCallbackServer.Server {
	return &server{
		handlers: handlers,
		sigfoxBackendDataCallbackMessageAdministrator: sigfoxBackendDataCallbackMessageAdministrator,
	}
}

func (s *server) HandleDataMessage(request *sigfoxBackendCallbackServer.HandleDataMessageRequest) (*sigfoxBackendCallbackServer.HandleDataMessageResponse, error) {
	for handlerIdx := range s.handlers {
		if s.handlers[handlerIdx].WantMessage(request.Message) {
			if err := s.handlers[handlerIdx].Handle(&sigfoxBackendDataMessageHandler.HandleRequest{
				Claims:      request.Claims,
				DataMessage: request.Message,
			}); err != nil {
				err = sigfoxBackendCallbackServerException.HandleDataMessage{Reasons: []string{"handling message", err.Error()}}
				log.Error(err.Error())
				return nil, err
			}
		}
	}
	return nil, nil
}

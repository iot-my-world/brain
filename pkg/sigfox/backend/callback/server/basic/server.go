package basic

import (
	"github.com/iot-my-world/brain/internal/log"
	sigfoxBackendDataMessageHandler "github.com/iot-my-world/brain/pkg/sigfox/backend/callback/data/message/handler"
	sigfoxBackendCallbackServer "github.com/iot-my-world/brain/pkg/sigfox/backend/callback/server"
	sigfoxBackendCallbackServerException "github.com/iot-my-world/brain/pkg/sigfox/backend/callback/server/exception"
)

type server struct {
	handlers []sigfoxBackendDataMessageHandler.Handler
}

func New(
	handlers []sigfoxBackendDataMessageHandler.Handler,
) sigfoxBackendCallbackServer.Server {
	return &server{
		handlers: handlers,
	}
}

func (s *server) MethodRequiresAuthorization(string) bool {
	return true
}

func (s *server) HandleDataMessage(request *sigfoxBackendCallbackServer.HandleDataMessageRequest) (*sigfoxBackendCallbackServer.HandleDataMessageResponse, error) {
	for handlerIdx := range s.handlers {
		wantsMessage, err := s.handlers[handlerIdx].WantMessage(request.Message)
		if err != nil {
			err = sigfoxBackendCallbackServerException.HandleDataMessage{Reasons: []string{"determining if handler wants message", err.Error()}}
			log.Error(err.Error())
			return nil, err
		}
		if wantsMessage {
			err := s.handlers[handlerIdx].Handle(request.Message)
			if err != nil {
				err = sigfoxBackendCallbackServerException.HandleDataMessage{Reasons: []string{"handling message", err.Error()}}
				log.Error(err.Error())
				return nil, err
			}
		}
	}
	return nil, nil
}

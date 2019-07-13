package basic

import (
	sigfoxBackendDataMessageHandler "github.com/iot-my-world/brain/pkg/sigfox/backend/callback/data/message/handler"
	sigfoxBackendCallbackServer "github.com/iot-my-world/brain/pkg/sigfox/backend/callback/server"
)

type server struct {
	handlers []sigfoxBackendDataMessageHandler.Handler
}

func New() sigfoxBackendCallbackServer.Server {
	return &server{}
}

func (s *server) MethodRequiresAuthorization(string) bool {
	return true
}

func (s *server) HandleDataMessage(*sigfoxBackendCallbackServer.HandleDataMessageRequest) (*sigfoxBackendCallbackServer.HandleDataMessageResponse, error) {
	return nil, nil
}

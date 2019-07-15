package message

import (
	sigfoxBackendDataDataCallbackMessage "github.com/iot-my-world/brain/pkg/sigfox/backend/callback/data/message"
	sigfoxBackendDataDataCallbackMessageHandler "github.com/iot-my-world/brain/pkg/sigfox/backend/callback/data/message/handler"
)

type handler struct {
}

func NewHandler() sigfoxBackendDataDataCallbackMessageHandler.Handler {
	return &handler{}
}

func (h *handler) WantMessage(dataMessage sigfoxBackendDataDataCallbackMessage.Message) bool {
	return true
}

func (h *handler) Handle(sigfoxBackendDataDataCallbackMessage.Message) error {
	return nil
}

package handler

import (
	sigfoxBackendDataDataCallbackMessage "github.com/iot-my-world/brain/pkg/sigfox/backend/callback/data/message"
)

type Handler interface {
	Handle(sigfoxBackendDataDataCallbackMessage.Message) error
	WantMessage(sigfoxBackendDataDataCallbackMessage.Message) bool
}

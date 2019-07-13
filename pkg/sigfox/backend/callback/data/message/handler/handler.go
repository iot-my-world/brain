package handler

import (
	sigfoxBackendDataCallbackMessage "github.com/iot-my-world/brain/pkg/sigfox/backend/callback/data/message"
)

type Handler interface {
	Handle(sigfoxBackendDataCallbackMessage.Message) error
	WantMessage(sigfoxBackendDataCallbackMessage.Message) (bool, error)
}

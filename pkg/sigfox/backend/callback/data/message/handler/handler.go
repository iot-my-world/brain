package handler

import (
	"github.com/iot-my-world/brain/pkg/security/claims"
	sigfoxBackendDataDataCallbackMessage "github.com/iot-my-world/brain/pkg/sigfox/backend/callback/data/message"
)

type Handler interface {
	Handle(*HandleRequest) error
	WantMessage(sigfoxBackendDataDataCallbackMessage.Message) bool
}

type HandleRequest struct {
	Claims      claims.Claims
	DataMessage sigfoxBackendDataDataCallbackMessage.Message
}

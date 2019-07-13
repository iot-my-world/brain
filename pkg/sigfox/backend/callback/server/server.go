package server

import (
	sigfoxBackendDataMessage "github.com/iot-my-world/brain/pkg/sigfox/backend/callback/data/message"
)

type Server interface {
	HandleDataMessage(*HandleDataMessageRequest) (*HandleDataMessageResponse, error)
}

const ServiceProvider = "SigfoxBackendCallbackServer"

const HandleDataMessageService = ServiceProvider + ".HandleDataMessage"

type HandleDataMessageRequest struct {
	Message sigfoxBackendDataMessage.Message
}

type HandleDataMessageResponse struct {
}

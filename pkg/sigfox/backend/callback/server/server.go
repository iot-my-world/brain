package server

import (
	sigfoxBackendDataMessage "github.com/iot-my-world/brain/pkg/sigfox/backend/callback/data/message"
)

type Server interface {
	HandleDataMessage(*HandleDataMessageRequest) (*HandleDataMessageResponse, error)
}

type HandleDataMessageRequest struct {
	Message sigfoxBackendDataMessage.Message
}

type HandleDataMessageResponse struct {
}
package server

import (
	"github.com/iot-my-world/brain/pkg/security/claims"
	sigfoxBackendDataMessage "github.com/iot-my-world/brain/pkg/sigfox/backend/callback/data/message"
)

type Server interface {
	HandleDataMessage(*HandleDataMessageRequest) (*HandleDataMessageResponse, error)
}

const ServiceProvider = "SigfoxBackendCallbackServer"

const HandleDataMessageService = ServiceProvider + ".HandleDataMessage"

type HandleDataMessageRequest struct {
	Claims  claims.Claims
	Message sigfoxBackendDataMessage.Message
}

type HandleDataMessageResponse struct {
}

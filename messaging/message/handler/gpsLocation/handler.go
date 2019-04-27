package gpsLocation

import (
	"gitlab.com/iotTracker/brain/log"
	messagingMessage "gitlab.com/iotTracker/messaging/message"
	messagingMessageHandler "gitlab.com/iotTracker/messaging/message/handler"
)

type handler struct {
}

func New() messagingMessageHandler.Handler {
	return &handler{}
}

func (h *handler) WantsMessage(message messagingMessage.Message) bool {
	return message.Type() == messagingMessage.GPSLocation
}

func (h *handler) HandleMessage(message messagingMessage.Message) {
	log.Info("handling gps message!!!!")
}

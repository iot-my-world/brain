package gps

import (
	brainException "gitlab.com/iotTracker/brain/exception"
	"gitlab.com/iotTracker/brain/log"
	messagingException "gitlab.com/iotTracker/messaging/exception"
	messagingMessage "gitlab.com/iotTracker/messaging/message"
	messagingMessageHandler "gitlab.com/iotTracker/messaging/message/handler"
	zx303GPSReadingMessage "gitlab.com/iotTracker/messaging/message/zx303/reading/gps"
)

type handler struct {
}

func New() messagingMessageHandler.Handler {
	return &handler{}
}

func (h *handler) WantsMessage(message messagingMessage.Message) bool {
	return message.Type() == messagingMessage.ZX303GPSReading
}

func (*handler) ValidateMessage(message messagingMessage.Message) error {
	reasonsInvalid := make([]string, 0)

	if _, ok := message.(zx303GPSReadingMessage.Message); !ok {
		reasonsInvalid = append(reasonsInvalid, "cannot cast message to zx303GPSReadingMessage.Message")
	}

	if len(reasonsInvalid) > 0 {
		return messagingException.InvalidMessage{Reasons: reasonsInvalid}
	}

	return nil
}

func (h *handler) HandleMessage(message messagingMessage.Message) error {
	if err := h.ValidateMessage(message); err != nil {
		return err
	}

	gpsReadingMessage, ok := message.(zx303GPSReadingMessage.Message)
	if !ok {
		return brainException.Unexpected{Reasons: []string{"cannot cast message to zx303GPSReadingMessage.Message"}}
	}

	log.Info("handling gps message!", gpsReadingMessage.Reading)

	return nil
}

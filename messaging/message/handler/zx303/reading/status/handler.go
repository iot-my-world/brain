package status

import (
	brainException "gitlab.com/iotTracker/brain/exception"
	humanUserLoginClaims "gitlab.com/iotTracker/brain/security/claims/login/user/human"
	zx303StatusReadingAdministrator "gitlab.com/iotTracker/brain/tracker/zx303/reading/status/administrator"
	messagingException "gitlab.com/iotTracker/messaging/exception"
	messagingMessage "gitlab.com/iotTracker/messaging/message"
	messagingMessageHandler "gitlab.com/iotTracker/messaging/message/handler"
	zx303StatusReadingMessage "gitlab.com/iotTracker/messaging/message/zx303/reading/status"
)

type handler struct {
	systemClaims                    *humanUserLoginClaims.Login
	zx303StatusReadingAdministrator zx303StatusReadingAdministrator.Administrator
}

func New(
	systemClaims *humanUserLoginClaims.Login,
	zx303StatusReadingAdministrator zx303StatusReadingAdministrator.Administrator,
) messagingMessageHandler.Handler {
	return &handler{
		systemClaims:                    systemClaims,
		zx303StatusReadingAdministrator: zx303StatusReadingAdministrator,
	}
}

func (h *handler) WantsMessage(message messagingMessage.Message) bool {
	return message.Type() == messagingMessage.ZX303StatusReading
}

func (*handler) ValidateMessage(message messagingMessage.Message) error {
	reasonsInvalid := make([]string, 0)

	if _, ok := message.(zx303StatusReadingMessage.Message); !ok {
		reasonsInvalid = append(reasonsInvalid, "cannot cast message to zx303StatusReadingMessage.Message")
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

	statusReadingMessage, ok := message.(zx303StatusReadingMessage.Message)
	if !ok {
		return brainException.Unexpected{Reasons: []string{"cannot cast message to zx303StatusReadingMessage.Message"}}
	}

	if _, err := h.zx303StatusReadingAdministrator.Create(&zx303StatusReadingAdministrator.CreateRequest{
		Claims:             h.systemClaims,
		ZX303StatusReading: statusReadingMessage.Reading,
	}); err != nil {
		return err
	}

	return nil
}

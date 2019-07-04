package status

import (
	brainException "github.com/iot-my-world/brain/internal/exception"
	zx303StatusReadingAdministrator "github.com/iot-my-world/brain/pkg/tracker/zx303/reading/status/administrator"
	humanUserLoginClaims "github.com/iot-my-world/brain/security/claims/login/user/human"
	messagingException "github.com/iot-my-world/messaging/exception"
	messagingMessage "github.com/iot-my-world/messaging/message"
	messagingMessageHandler "github.com/iot-my-world/messaging/message/handler"
	zx303StatusReadingMessage "github.com/iot-my-world/messaging/message/zx303/reading/status"
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

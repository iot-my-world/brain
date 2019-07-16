package handler

import (
	"encoding/binary"
	"fmt"
	"github.com/iot-my-world/brain/pkg/device/sigbug/sigfox/message"
	sigfoxBackendDataDataCallbackMessageHandlerException "github.com/iot-my-world/brain/pkg/device/sigbug/sigfox/message/handler/exception"
	sigfoxBackendDataDataCallbackMessage "github.com/iot-my-world/brain/pkg/sigfox/backend/callback/data/message"
	sigfoxBackendDataDataCallbackMessageHandler "github.com/iot-my-world/brain/pkg/sigfox/backend/callback/data/message/handler"
	"math"
)

type handler struct {
}

func New() sigfoxBackendDataDataCallbackMessageHandler.Handler {
	return &handler{}
}

func (h *handler) WantMessage(dataMessage sigfoxBackendDataDataCallbackMessage.Message) bool {
	if len(dataMessage.Data) == 0 {
		return false
	}

	switch dataMessage.Data[0] {
	case message.GPSReading:
		if len(dataMessage.Data) == 9 {
			return true
		}
	}

	return false
}

func (h *handler) Handle(dataMessage sigfoxBackendDataDataCallbackMessage.Message) error {

	switch dataMessage.Data[0] {
	case message.GPSReading:
		return h.handleGPSMessage(dataMessage)
	}

	return nil
}

func (h *handler) handleGPSMessage(dataMessage sigfoxBackendDataDataCallbackMessage.Message) error {
	if len(dataMessage.Data) != 9 {
		return sigfoxBackendDataDataCallbackMessageHandlerException.HandleGPSMessage{Reasons: []string{"message data not long enough"}}
	}

	latitude := math.Float32frombits(binary.LittleEndian.Uint32(dataMessage.Data[1:5]))
	longitude := math.Float32frombits(binary.LittleEndian.Uint32(dataMessage.Data[5:]))

	fmt.Printf("%f, %f", latitude, longitude)

	return nil
}

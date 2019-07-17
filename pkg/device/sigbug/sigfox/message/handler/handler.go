package handler

import (
	"encoding/binary"
	"fmt"
	sigbugRecordHandler "github.com/iot-my-world/brain/pkg/device/sigbug/recordHandler"
	"github.com/iot-my-world/brain/pkg/device/sigbug/sigfox/message"
	sigfoxBackendDataDataCallbackMessageHandlerException "github.com/iot-my-world/brain/pkg/device/sigbug/sigfox/message/handler/exception"
	sigfoxBackendDataDataCallbackMessage "github.com/iot-my-world/brain/pkg/sigfox/backend/callback/data/message"
	sigfoxBackendDataDataCallbackMessageHandler "github.com/iot-my-world/brain/pkg/sigfox/backend/callback/data/message/handler"
	"math"
)

type handler struct {
	sigbugRecordHandler sigbugRecordHandler.RecordHandler
}

func New(
	sigbugRecordHandler sigbugRecordHandler.RecordHandler,
) sigfoxBackendDataDataCallbackMessageHandler.Handler {
	return &handler{
		sigbugRecordHandler: sigbugRecordHandler,
	}
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

func (h *handler) Handle(request *sigfoxBackendDataDataCallbackMessageHandler.HandleRequest) error {

	switch request.DataMessage.Data[0] {
	case message.GPSReading:
		return h.handleGPSMessage(request)
	}

	return nil
}

func (h *handler) handleGPSMessage(request *sigfoxBackendDataDataCallbackMessageHandler.HandleRequest) error {
	if len(request.DataMessage.Data) != 9 {
		return sigfoxBackendDataDataCallbackMessageHandlerException.HandleGPSMessage{Reasons: []string{"message data not long enough"}}
	}

	latitude := math.Float32frombits(binary.LittleEndian.Uint32(request.DataMessage.Data[1:5]))
	longitude := math.Float32frombits(binary.LittleEndian.Uint32(request.DataMessage.Data[5:]))

	fmt.Printf("%f, %f", latitude, longitude)

	return nil
}

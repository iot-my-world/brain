package message

import "github.com/iot-my-world/brain/pkg/search/identifier/id"

type Message struct {
	DeviceIdentifier id.Identifier
	Data             []byte
}

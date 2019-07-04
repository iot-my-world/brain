package websocket

import "github.com/iot-my-world/brain/pkg/communication/websocket/message"

type Message struct {
	Type        message.ExoWSMsgType `json:"type"`
	SerialData  string               `json:"serialData"`
	ReBroadcast bool                 `json:"reBroadcast"`
}

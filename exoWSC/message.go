package exoWSC

import "github.com/iot-my-world/brain/exoWSC/message"

type Message struct {
	Type        message.ExoWSMsgType `json:"type"`
	SerialData  string               `json:"serialData"`
	ReBroadcast bool                 `json:"reBroadcast"`
}

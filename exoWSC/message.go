package exoWSC

import "gitlab.com/iotTracker/brain/exoWSC/message"

type Message struct {
	Type        message.ExoWSMsgType `json:"type"`
	SerialData  string               `json:"serialData"`
	ReBroadcast bool                 `json:"reBroadcast"`
}

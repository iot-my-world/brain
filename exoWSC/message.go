package exoWSC

import "bitbucket.org/gotimekeeper/exoWSC/message"

type Message struct {
	Type        message.ExoWSMsgType `json:"type"`
	SerialData  string               `json:"serialData"`
	ReBroadcast bool                 `json:"reBroadcast"`
}

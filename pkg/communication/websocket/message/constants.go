package message

type ExoWSMsgType int

const (
	CreateServiceContextRequest ExoWSMsgType = iota
	CreateServiceContextResponse
	GetServiceContextRequest
	GetServiceContextResponse
	ClockEvent
	WelcomeMessage
)

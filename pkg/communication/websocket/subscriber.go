package websocket

type Subscriber interface {
	Send(message Message) error
}

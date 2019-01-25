package exoWSC

type Subscriber interface {
	Send(message Message) error
}

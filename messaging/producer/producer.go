package producer

type Producer interface {
	Start() error
	Produce(data []byte) error
}

package message

type Message struct {
	DeviceId string `json:"backendId"`
	Data     []byte `json:"data"`
}

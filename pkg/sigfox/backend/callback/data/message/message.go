package message

type Message struct {
	DeviceId string `json:"deviceId"`
	Data     []byte `json:"data"`
}

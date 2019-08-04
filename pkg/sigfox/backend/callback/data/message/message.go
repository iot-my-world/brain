package message

type Message struct {
	Id       string `json:"id" bson:"id"`
	DeviceId string `json:"deviceId" bson:"deviceId"`
	Data     []byte `json:"data" bson:"data"`
}

func (m *Message) SetId(id string) {
	m.Id = id
}

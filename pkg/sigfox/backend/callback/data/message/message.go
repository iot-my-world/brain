package message

type Message struct {
	Id        string `json:"id" bson:"id"`
	Timestamp int    `json:"timeStamp" bson:"timeStamp"`
	DeviceId  string `json:"deviceId" bson:"deviceId"`
	Data      []byte `json:"data" bson:"data"`
}

func (m *Message) SetId(id string) {
	m.Id = id
}

package rfId

type TagEvent struct {
	TagId string `json:"tag_id" bson:"tag_id"`
	TagTime int64 `json:"tag_time" bson:"tag_time"`
	ReceivedTime int64 `json:"received_time" bson:"received_time"`
}

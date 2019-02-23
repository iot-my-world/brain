package reading

type Reading struct {
	TimeStamp int64   `json:"timeStamp" bson:"timeStamp"`
	Id        string  `json:"id" bson:"id"`
	IMEI      string  `json:"imei" bson:"imei"`
	Raw       string  `json:"raw" bson:"raw"`
	Latitude  float32 `json:"latitude" bson:"latitude"`
	Longitude float32 `json:"longitude" bson:"longitude"`
}
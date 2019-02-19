package reading

type Reading struct {
	Id              string `json:"id" bson:"id"`
	IMEI            string `json:"imei" bson:"imei"`
	Raw             string `json:"raw" bson:"raw"`
	SouthCoordinate string `json:"southCoordinate" bson:"southCoordinate"`
	EastCoordinate  string `json:"eastCoordinate" bson:"eastCoordinate"`
	TimeStamp       int64  `json:"timeStamp" bson:"timeStamp"`
}

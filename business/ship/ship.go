package ship

type Ship struct{
	Id string `json:"id" bson:"id"`
	Name string `json:"name" bson:"name"`
	Berth string `json:"berth" bson:"berth"`
	InDateTime int64 `json:"inDateTime" bson:"inDateTime"`
	OutDateTime int64 `json:"outDateTime" bson:"outDateTime"`
	Deleted bool `json:"deleted" bson:"deleted"`
}

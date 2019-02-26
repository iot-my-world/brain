package system

type System struct {
	Id                string `json:"id" bson:"id"`
	Name              string `json:"name" bson:"name"`
	AdminEmailAddress string `json:"adminEmailAddress" bson:"adminEmailAddress"`
}

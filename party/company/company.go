package company

type Company struct {
	Id string `json:"id" bson:"id"`
	Name string `json:"name" bson:"name"`
}

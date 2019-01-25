package businessRole

type BusinessRole struct{
	Id string `json:"id" bson:"id"`
	Name string `json:"name" bson:"name"`
	PayRate float32 `json:"payRate" bson:"payRate"`
}
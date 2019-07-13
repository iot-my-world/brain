package backend

type Backend struct {
	Id    string `json:"id" bson:"id"`
	Name  string `json:"name" bson:"name"`
	Token string `json:"token" bson:"token"`
}

func (b Backend) SetId(id string) {
	b.Id = id
}

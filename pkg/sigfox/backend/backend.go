package backend

type Backend struct {
	Id   string `json:"id" bson:"id"`
	Name string `json:"name" bson:"name"`
}

func (b Backend) SetId(id string) {
	b.Id = id
}

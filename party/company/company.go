package company

type Company struct {
	Id   string `json:"id" bson:"id"`
	Name string `json:"name" bson:"name"`
	// The email address which will be used to invite the admin
	// user of the company
	// I.e. the first user of the system from the company
	AdminEmail string `json:"adminEmail" bson:"adminEmail"`
}

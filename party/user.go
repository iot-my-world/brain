package party

// Defines the User record for the database
type User struct {
	Id string `json:"id" bson:"id"`

	// Personal Details
	Name    string `json:"name" bson:"name"`
	Surname string `json:"surname" bson:"surname"`

	// System Details
	Username     string `json:"username" bson:"username"`
	EmailAddress string `json:"emailAddress" bson:"emailAddress"`
	Password     []byte `json:"pwd" bson:"pwd"`
	SystemRole   string `json:"systemRole" bson:"systemRole"`
}

type NewUser struct {
	// Personal Details
	Name    string `json:"name"`
	Surname string `json:"surname"`
	IDNo    int    `json:"idNo"`

	// System Details
	Username     string `json:"username"`
	EmailAddress string `json:"emailAddress"`
	Password     string `json:"password"`
	SystemRole   string `json:"systemRole"`
}
package user

// Defines the User record for the database
type User struct {
	// Personal Details
	// TODO: Split out into "PersonalDetails" Struct
	Name    string `json:"name" bson:"name"`
	Surname string `json:"surname" bson:"surname"`
	IDNo    int    `json:"idNo" bson:"idNo"`

	// System Details
	Username   string `json:"username" bson:"username"`
	Password   []byte `json:"pwd" bson:"pwd"`
	SystemRole string `json:"systemRole" bson:"systemRole"`
	TagID      string `json:"tagID" bson:"tagID"`

	// Business Details
	// TODO: Split out into "BusinessDetails" Struct
	BusinessRole string `json:"businessRole" bson:"businessRole"`
}

type NewUser struct {
	// Personal Details
	// TODO: Split out into "PersonalDetails" Struct
	Name    string `json:"name"`
	Surname string `json:"surname"`
	IDNo    int    `json:"idNo"`

	// System Details
	Username   string `json:"username"`
	Password   string `json:"password"`
	SystemRole string `json:"systemRole"`
	TagID      string `json:"tagID"`

	// Business Details
	// TODO: Split out into "BusinessDetails" Struct
	BusinessRole string `json:"businessRole"`
}
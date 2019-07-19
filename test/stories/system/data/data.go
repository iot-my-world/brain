package data

import (
	"github.com/iot-my-world/brain/pkg/sigfox/backend"
	humanUser "github.com/iot-my-world/brain/pkg/user/human"
	sigfoxBackendTestModule "github.com/iot-my-world/brain/test/modules/sigfox/backend"
)

var User = humanUser.User{
	//Id string

	// Personal Details
	Name:    "root",
	Surname: "root",

	// System Details
	Username: "root",
	// EmailAddress
	Password: []byte("12345"),
	Roles:    make([]string, 0),

	// Party Details
	//ParentPartyType
	//ParentId
	//PartyType
	//PartyId
}

var SigfoxBackendTestData = []sigfoxBackendTestModule.Data{
	{
		Backend: backend.Backend{
			Name: "Sigfox",
		},
	},
}

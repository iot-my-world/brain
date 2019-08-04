package data

import (
	"github.com/iot-my-world/brain/pkg/device/sigbug"
	"github.com/iot-my-world/brain/pkg/party"
	"github.com/iot-my-world/brain/pkg/party/client"
	"github.com/iot-my-world/brain/pkg/search/identifier/id"
	sigfoxBackendDataCallbackMessage "github.com/iot-my-world/brain/pkg/sigfox/backend/callback/data/message"
	humanUser "github.com/iot-my-world/brain/pkg/user/human"
	clientTest "github.com/iot-my-world/brain/test/modules/party/client"
)

var TestData = map[string][]struct {
	ClientTestData clientTest.Data
	SigbugDevices  []sigbug.Sigbug
}{
	"root": {
		{
			ClientTestData: clientTest.Data{
				Client: client.Client{
					Name:              "John",
					Type:              client.Individual,
					AdminEmailAddress: "john@gmail.com",
					// ParentId: // populated on creation
				},
				AdminUser: humanUser.User{
					Name:     "John",
					Surname:  "Smith",
					Username: "jsmith",
					Password: []byte("123"),
					Roles:    make([]string, 0),
				},
				Users: []humanUser.User{
					{
						Name:         "Samantha",
						Surname:      "Smith",
						Username:     "ssmith",
						Password:     []byte("123"),
						EmailAddress: "sam@gmail.com",
						Roles:        make([]string, 0),
					},
				},
			},
			SigbugDevices: []sigbug.Sigbug{
				{
					DeviceId:       "1asdfepoi",
					OwnerPartyType: party.Client,
					OwnerId: id.Identifier{
						Id: "John", // populated correctly on creation
					},
					AssignedPartyType: party.Client,
					AssignedId: id.Identifier{
						Id: "John", // populated correctly on creation
					},
					LastMessage: sigfoxBackendDataCallbackMessage.Message{
						Data: make([]byte, 0),
					},
				},
				{
					DeviceId:       "1aspoive",
					OwnerPartyType: party.Client,
					OwnerId: id.Identifier{
						Id: "John", // populated correctly on creation
					},
					AssignedPartyType: party.Client,
					AssignedId: id.Identifier{
						Id: "John", // populated correctly on creation
					},
					LastMessage: sigfoxBackendDataCallbackMessage.Message{
						Data: make([]byte, 0),
					},
				},
			},
		},
		{
			ClientTestData: clientTest.Data{
				Client: client.Client{
					Name:              "Sunbeam Tech",
					Type:              client.Company,
					AdminEmailAddress: "sunbeamTechAdmin@sunbeam.com",
					// ParentId: // populated on creation
				},
				AdminUser: humanUser.User{
					Name:     "Johan",
					Surname:  "von delft",
					Username: "sunbeamAdmin",
					Password: []byte("123"),
					Roles:    make([]string, 0),
				},
				Users: []humanUser.User{
					{
						Name:         "Juliana",
						Surname:      "Trump",
						Username:     "sunbeamUser1",
						Password:     []byte("123"),
						EmailAddress: "sunbeamUser1@sunbeam.com",
						Roles:        make([]string, 0),
					},
					{
						Name:         "Paul",
						Surname:      "Xulu",
						Username:     "sunbeamUser2",
						Password:     []byte("123"),
						EmailAddress: "sunbeamUser2@sunbeam.com",
						Roles:        make([]string, 0),
					},
				},
			},
			SigbugDevices: []sigbug.Sigbug{
				{
					DeviceId:       "2aposdiu3nv",
					OwnerPartyType: party.Client,
					OwnerId: id.Identifier{
						Id: "Sunbeam Tech",
					},
					AssignedPartyType: party.Client,
					AssignedId: id.Identifier{
						Id: "Sunbeam Tech", // populated correctly on creation
					},
					LastMessage: sigfoxBackendDataCallbackMessage.Message{
						Data: make([]byte, 0),
					},
				},
				{
					DeviceId:       "2piaseinve",
					OwnerPartyType: party.Client,
					OwnerId: id.Identifier{
						Id: "Sunbeam Tech", // populated correctly on creation
					},
					AssignedPartyType: party.Client,
					AssignedId: id.Identifier{
						Id: "Sunbeam Tech", // populated correctly on creation
					},
					LastMessage: sigfoxBackendDataCallbackMessage.Message{
						Data: make([]byte, 0),
					},
				},
			},
		},
	},

	"Monteagle Logistics Limited": {
		{
			ClientTestData: clientTest.Data{
				Client: client.Client{
					Name:              "Pick 'n Pay",
					Type:              client.Company,
					AdminEmailAddress: "picknpayAdmin@picknpay.com",
					// ParentId: // populated on creation
				},
				AdminUser: humanUser.User{
					Name:     "Yoland",
					Surname:  "Govender",
					Username: "picknpayAdmin",
					Password: []byte("123"),
					Roles:    make([]string, 0),
				},
				Users: []humanUser.User{
					{
						Name:         "picknpayUser1",
						Surname:      "test1",
						Username:     "picknpayUser1",
						Password:     []byte("123"),
						EmailAddress: "picknpayUser1@picknpay.com",
						Roles:        make([]string, 0),
					},
					{
						Name:         "picknpayUser2",
						Surname:      "test2",
						Username:     "picknpayUser2",
						Password:     []byte("123"),
						EmailAddress: "picknpayUser2@picknpay.com",
						Roles:        make([]string, 0),
					},
				},
			},
			SigbugDevices: []sigbug.Sigbug{
				{
					DeviceId:       "3poisdoiece",
					OwnerPartyType: party.Client,
					OwnerId: id.Identifier{
						Id: "Pick 'n Pay", // populated correctly on creation
					},
					AssignedPartyType: party.Client,
					AssignedId: id.Identifier{
						Id: "Pick 'n Pay", // populated correctly on creation
					},
					LastMessage: sigfoxBackendDataCallbackMessage.Message{
						Data: make([]byte, 0),
					},
				},
				{
					DeviceId:       "3poiasdvpoie",
					OwnerPartyType: party.Client,
					OwnerId: id.Identifier{
						Id: "Pick 'n Pay", // populated correctly on creation
					},
					AssignedPartyType: party.Client,
					AssignedId: id.Identifier{
						Id: "Pick 'n Pay", // populated correctly on creation
					},
					LastMessage: sigfoxBackendDataCallbackMessage.Message{
						Data: make([]byte, 0),
					},
				},
			},
		},

		{
			ClientTestData: clientTest.Data{
				Client: client.Client{
					Name:              "Woolworths",
					Type:              client.Company,
					AdminEmailAddress: "woolworthsAdmin@woolworths.com",
					// ParentId: // populated on creation
				},
				AdminUser: humanUser.User{
					Name:     "Pauline",
					Surname:  "Kruger",
					Username: "woolworthsAdmin",
					Password: []byte("123"),
					Roles:    make([]string, 0),
				},
				Users: []humanUser.User{
					{
						Name:         "woolworthsUser1",
						Surname:      "test1",
						Username:     "woolworthsUser1",
						Password:     []byte("123"),
						EmailAddress: "woolworthsUser1@woolworths.com",
						Roles:        make([]string, 0),
					},
					{
						Name:         "woolworthsUser2",
						Surname:      "test2",
						Username:     "woolworthsUser2",
						Password:     []byte("123"),
						EmailAddress: "woolworthsUser2@woolworths.com",
						Roles:        make([]string, 0),
					},
				},
			},
			SigbugDevices: []sigbug.Sigbug{
				{
					DeviceId:       "4paspoivee",
					OwnerPartyType: party.Client,
					OwnerId: id.Identifier{
						Id: "Woolworths", // populated correctly on creation
					},
					AssignedPartyType: party.Client,
					AssignedId: id.Identifier{
						Id: "Woolworths", // populated correctly on creation
					},
					LastMessage: sigfoxBackendDataCallbackMessage.Message{
						Data: make([]byte, 0),
					},
				},
				{
					DeviceId:       "4poiadoiee",
					OwnerPartyType: party.Client,
					OwnerId: id.Identifier{
						Id: "Woolworths", // populated correctly on creation
					},
					AssignedPartyType: party.Client,
					AssignedId: id.Identifier{
						Id: "Woolworths", // populated correctly on creation
					},
					LastMessage: sigfoxBackendDataCallbackMessage.Message{
						Data: make([]byte, 0),
					},
				},
			},
		},
	},

	"DHL": {
		{
			ClientTestData: clientTest.Data{
				Client: client.Client{
					Name:              "Makro",
					Type:              client.Company,
					AdminEmailAddress: "makroAdmin@makro.com",
					// ParentId: // populated on creation
				},
				AdminUser: humanUser.User{
					Name:     "Johan",
					Surname:  "Smith",
					Username: "makroAdmin",
					Password: []byte("123"),
					Roles:    make([]string, 0),
				},
				Users: []humanUser.User{
					{
						Name:         "makroUser1",
						Surname:      "test1",
						Username:     "makroUser1",
						Password:     []byte("123"),
						EmailAddress: "makroUser1@makro.com",
						Roles:        make([]string, 0),
					},
					{
						Name:         "makroUser2",
						Surname:      "test2",
						Username:     "makroUser2",
						Password:     []byte("123"),
						EmailAddress: "makroUser2@makro.com",
						Roles:        make([]string, 0),
					},
				},
			},
			SigbugDevices: []sigbug.Sigbug{
				{
					DeviceId:       "5p0aseojve",
					OwnerPartyType: party.Client,
					OwnerId: id.Identifier{
						Id: "Makro", // populated correctly on creation
					},
					AssignedPartyType: party.Client,
					AssignedId: id.Identifier{
						Id: "Makro", // populated correctly on creation
					},
					LastMessage: sigfoxBackendDataCallbackMessage.Message{
						Data: make([]byte, 0),
					},
				},
				{
					DeviceId:       "5poieoiae",
					OwnerPartyType: party.Client,
					OwnerId: id.Identifier{
						Id: "Makro", // populated correctly on creation
					},
					AssignedPartyType: party.Client,
					AssignedId: id.Identifier{
						Id: "Makro", // populated correctly on creation
					},
					LastMessage: sigfoxBackendDataCallbackMessage.Message{
						Data: make([]byte, 0),
					},
				},
			},
		},

		{
			ClientTestData: clientTest.Data{
				Client: client.Client{
					Name:              "Fruit 'n Veg",
					Type:              client.Company,
					AdminEmailAddress: "fruitnvegAdmin@fruitnveg.com",
					// ParentId: // populated on creation
				},
				AdminUser: humanUser.User{
					Name:     "Vicky",
					Surname:  "smith",
					Username: "fruitnvegAdmin",
					Password: []byte("123"),
					Roles:    make([]string, 0),
				},
				Users: []humanUser.User{
					{
						Name:         "fruitnvegUser1",
						Surname:      "test1",
						Username:     "fruitnvegUser1",
						Password:     []byte("123"),
						EmailAddress: "fruitnvegUser1@fruitnveg.com",
						Roles:        make([]string, 0),
					},
					{
						Name:         "fruitnvegUser2",
						Surname:      "test2",
						Username:     "fruitnvegUser2",
						Password:     []byte("123"),
						EmailAddress: "fruitnvegUser2@fruitnveg.com",
						Roles:        make([]string, 0),
					},
				},
			},
			SigbugDevices: []sigbug.Sigbug{
				{
					DeviceId:       "6eiuoipoi",
					OwnerPartyType: party.Client,
					OwnerId: id.Identifier{
						Id: "Fruit 'n Veg", // populated correctly on creation
					},
					AssignedPartyType: party.Client,
					AssignedId: id.Identifier{
						Id: "Fruit 'n Veg", // populated correctly on creation
					},
					LastMessage: sigfoxBackendDataCallbackMessage.Message{
						Data: make([]byte, 0),
					},
				},
				{
					DeviceId:       "6poijwepoi",
					OwnerPartyType: party.Client,
					OwnerId: id.Identifier{
						Id: "Fruit 'n Veg", // populated correctly on creation
					},
					AssignedPartyType: party.Client,
					AssignedId: id.Identifier{
						Id: "Fruit 'n Veg", // populated correctly on creation
					},
					LastMessage: sigfoxBackendDataCallbackMessage.Message{
						Data: make([]byte, 0),
					},
				},
			},
		},
	},

	"Reinhard Trucking": {
		{
			ClientTestData: clientTest.Data{
				Client: client.Client{
					Name:              "Spar",
					Type:              client.Company,
					AdminEmailAddress: "sparAdmin@spar.com",
					// ParentId: // populated on creation
				},
				AdminUser: humanUser.User{
					Name:     "Sipho",
					Surname:  "Shezi",
					Username: "sparAdmin",
					Password: []byte("123"),
					Roles:    make([]string, 0),
				},
				Users: []humanUser.User{
					{
						Name:         "sparUser1",
						Surname:      "test1",
						Username:     "sparUser1",
						Password:     []byte("123"),
						EmailAddress: "sparUser1@spar.com",
						Roles:        make([]string, 0),
					},
					{
						Name:         "sparUser2",
						Surname:      "test2",
						Username:     "sparUser2",
						Password:     []byte("123"),
						EmailAddress: "sparUser2@spar.com",
						Roles:        make([]string, 0),
					},
				},
			},
			SigbugDevices: []sigbug.Sigbug{
				{
					DeviceId:       "7pasoiepoije",
					OwnerPartyType: party.Client,
					OwnerId: id.Identifier{
						Id: "Spar", // populated correctly on creation
					},
					AssignedPartyType: party.Client,
					AssignedId: id.Identifier{
						Id: "Spar", // populated correctly on creation
					},
					LastMessage: sigfoxBackendDataCallbackMessage.Message{
						Data: make([]byte, 0),
					},
				},
				{
					DeviceId:       "7aoiweoijvepi",
					OwnerPartyType: party.Client,
					OwnerId: id.Identifier{
						Id: "Spar", // populated correctly on creation
					},
					AssignedPartyType: party.Client,
					AssignedId: id.Identifier{
						Id: "Spar", // populated correctly on creation
					},
					LastMessage: sigfoxBackendDataCallbackMessage.Message{
						Data: make([]byte, 0),
					},
				},
			},
		},

		{
			ClientTestData: clientTest.Data{
				Client: client.Client{
					Name:              "Game",
					Type:              client.Company,
					AdminEmailAddress: "gameAdmin@game.com",
					// ParentId: // populated on creation
				},
				AdminUser: humanUser.User{
					Name:     "Michael",
					Surname:  "Black",
					Username: "gameAdmin",
					Password: []byte("123"),
					Roles:    make([]string, 0),
				},
				Users: []humanUser.User{
					{
						Name:         "gameUser1",
						Surname:      "test1",
						Username:     "gameUser1",
						Password:     []byte("123"),
						EmailAddress: "gameUser1@game.com",
						Roles:        make([]string, 0),
					},
					{
						Name:         "gameUser2",
						Surname:      "test2",
						Username:     "gameUser2",
						Password:     []byte("123"),
						EmailAddress: "gameUser2@game.com",
						Roles:        make([]string, 0),
					},
				},
			},
			SigbugDevices: []sigbug.Sigbug{
				{
					DeviceId:       "8aspoeipe",
					OwnerPartyType: party.Client,
					OwnerId: id.Identifier{
						Id: "Game", // populated correctly on creation
					},
					AssignedPartyType: party.Client,
					AssignedId: id.Identifier{
						Id: "Game", // populated correctly on creation
					},
					LastMessage: sigfoxBackendDataCallbackMessage.Message{
						Data: make([]byte, 0),
					},
				},
				{
					DeviceId:       "8iasepieoij",
					OwnerPartyType: party.Client,
					OwnerId: id.Identifier{
						Id: "Game", // populated correctly on creation
					},
					AssignedPartyType: party.Client,
					AssignedId: id.Identifier{
						Id: "Game", // populated correctly on creation
					},
					LastMessage: sigfoxBackendDataCallbackMessage.Message{
						Data: make([]byte, 0),
					},
				},
			},
		},
	},
}

package client

import (
	"github.com/iot-my-world/brain/party/client"
	clientTest "github.com/iot-my-world/brain/test/modules/party/client"
	humanUser "github.com/iot-my-world/brain/user/human"
)

var TestData = map[string][]clientTest.Data{
	"root": { // clients created by root
		{
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
			},
			Users: []humanUser.User{
				{
					Name:         "Samantha",
					Surname:      "Smith",
					Username:     "ssmith",
					Password:     []byte("123"),
					EmailAddress: "sam@gmail.com",
				},
			},
		},
		{
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
			},
			Users: []humanUser.User{
				{
					Name:         "Juliana",
					Surname:      "Trump",
					Username:     "sunbeamUser1",
					Password:     []byte("123"),
					EmailAddress: "sunbeamUser1@sunbeam.com",
				},
				{
					Name:         "Paul",
					Surname:      "Xulu",
					Username:     "sunbeamUser2",
					Password:     []byte("123"),
					EmailAddress: "sunbeamUser2@sunbeam.com",
				},
			},
		},
	},
	"Monteagle Logistics Limited": {
		{
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
			},
			Users: []humanUser.User{
				{
					Name:         "picknpayUser1",
					Surname:      "test1",
					Username:     "picknpayUser1",
					Password:     []byte("123"),
					EmailAddress: "picknpayUser1@picknpay.com",
				},
				{
					Name:         "picknpayUser2",
					Surname:      "test2",
					Username:     "picknpayUser2",
					Password:     []byte("123"),
					EmailAddress: "picknpayUser2@picknpay.com",
				},
			},
		},
		{
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
			},
			Users: []humanUser.User{
				{
					Name:         "woolworthsUser1",
					Surname:      "test1",
					Username:     "woolworthsUser1",
					Password:     []byte("123"),
					EmailAddress: "woolworthsUser1@woolworths.com",
				},
				{
					Name:         "woolworthsUser2",
					Surname:      "test2",
					Username:     "woolworthsUser2",
					Password:     []byte("123"),
					EmailAddress: "woolworthsUser2@woolworths.com",
				},
			},
		},
	},
	"DHL": {
		{
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
			},
			Users: []humanUser.User{
				{
					Name:         "makroUser1",
					Surname:      "test1",
					Username:     "makroUser1",
					Password:     []byte("123"),
					EmailAddress: "makroUser1@makro.com",
				},
				{
					Name:         "makroUser2",
					Surname:      "test2",
					Username:     "makroUser2",
					Password:     []byte("123"),
					EmailAddress: "makroUser2@makro.com",
				},
			},
		},
		{
			Client: client.Client{
				Name:              "Fruit n Veg",
				Type:              client.Company,
				AdminEmailAddress: "fruitnvegAdmin@fruitnveg.com",
				// ParentId: // populated on creation
			},
			AdminUser: humanUser.User{
				Name:     "Vicky",
				Surname:  "smith",
				Username: "fruitnvegAdmin",
				Password: []byte("123"),
			},
			Users: []humanUser.User{
				{
					Name:         "fruitnvegUser1",
					Surname:      "test1",
					Username:     "fruitnvegUser1",
					Password:     []byte("123"),
					EmailAddress: "fruitnvegUser1@fruitnveg.com",
				},
				{
					Name:         "fruitnvegUser2",
					Surname:      "test2",
					Username:     "fruitnvegUser2",
					Password:     []byte("123"),
					EmailAddress: "fruitnvegUser2@fruitnveg.com",
				},
			},
		},
	},
	"Reinhard Trucking": {
		{
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
			},
			Users: []humanUser.User{
				{
					Name:         "sparUser1",
					Surname:      "test1",
					Username:     "sparUser1",
					Password:     []byte("123"),
					EmailAddress: "sparUser1@spar.com",
				},
				{
					Name:         "sparUser2",
					Surname:      "test2",
					Username:     "sparUser2",
					Password:     []byte("123"),
					EmailAddress: "sparUser2@spar.com",
				},
			},
		},
		{
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
			},
			Users: []humanUser.User{
				{
					Name:         "gameUser1",
					Surname:      "test1",
					Username:     "gameUser1",
					Password:     []byte("123"),
					EmailAddress: "gameUser1@game.com",
				},
				{
					Name:         "gameUser2",
					Surname:      "test2",
					Username:     "gameUser2",
					Password:     []byte("123"),
					EmailAddress: "gameUser2@game.com",
				},
			},
		},
	},
}

package data

import (
	"github.com/iot-my-world/brain/pkg/party/client"
	humanUser "github.com/iot-my-world/brain/user/human"
)

type TestData struct {
	Client    client.Client
	AdminUser humanUser.User
	Users     []humanUser.User
}

var EntitiesAndAdminUsersToCreate = map[string][]TestData{
	"Monteagle Logistics Limited": {
		{
			Client: client.Client{
				Name:              "Pick 'n Pay",
				AdminEmailAddress: "picknpayAdmin@picknpay.com",
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
				AdminEmailAddress: "woolworthsAdmin@woolworths.com",
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
				AdminEmailAddress: "makroAdmin@makro.com",
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
				AdminEmailAddress: "fruitnvegAdmin@fruitnveg.com",
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
				AdminEmailAddress: "sparAdmin@spar.com",
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
				AdminEmailAddress: "gameAdmin@game.com",
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

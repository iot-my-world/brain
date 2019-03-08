package client

import (
	"gitlab.com/iotTracker/brain/party/client"
	"gitlab.com/iotTracker/brain/party/user"
)

type TestData struct {
	Client    client.Client
	AdminUser user.User
	Users     []user.User
}

var EntitiesAndAdminUsersToCreate = map[string][]TestData{
	"Monteagle Logistics Limited": {
		{
			Client: client.Client{
				Name:              "Pick 'n Pay",
				AdminEmailAddress: "admin@picknpay.com",
			},
			AdminUser: user.User{
				Name:     "Yoland",
				Surname:  "Govender",
				Username: "picknpayAdmin",
				Password: []byte("123"),
			},
			Users: []user.User{
				{
					Name:         "picknpayUser1",
					Surname:      "test1",
					Username:     "picknpayTestUser1",
					Password:     []byte("123"),
					EmailAddress: "picknpayUser1@picknpay.com",
				},
				{
					Name:         "picknpayUser2",
					Surname:      "test2",
					Username:     "picknpayTestUser2",
					Password:     []byte("123"),
					EmailAddress: "picknpayUser2@picknpay.com",
				},
			},
		},
		{
			Client: client.Client{
				Name:              "Woolworths",
				AdminEmailAddress: "admin@woolworths.com",
			},
			AdminUser: user.User{
				Name:     "Pauline",
				Surname:  "Kruger",
				Username: "woolworthsAdmin",
				Password: []byte("123"),
			},
			Users: []user.User{
				{
					Name:         "woolworthsUser1",
					Surname:      "test1",
					Username:     "woolworthsTestUser1",
					Password:     []byte("123"),
					EmailAddress: "woolworthsUser1@woolworths.com",
				},
				{
					Name:         "woolworthsUser2",
					Surname:      "test2",
					Username:     "woolworthsTestUser2",
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
				AdminEmailAddress: "admin@makro.com",
			},
			AdminUser: user.User{
				Name:     "Johan",
				Surname:  "Smith",
				Username: "makroAdmin",
				Password: []byte("123"),
			},
			Users: []user.User{
				{
					Name:         "makroUser1",
					Surname:      "test1",
					Username:     "makroTestUser1",
					Password:     []byte("123"),
					EmailAddress: "makroUser1@makro.com",
				},
				{
					Name:         "makroUser2",
					Surname:      "test2",
					Username:     "makroTestUser2",
					Password:     []byte("123"),
					EmailAddress: "woolworthsUser2@woolworths.com",
				},
			},
		},
		{
			Client: client.Client{
				Name:              "Fruit n Veg",
				AdminEmailAddress: "admin@fruitnveg.com",
			},
			AdminUser: user.User{
				Name:     "Vicky",
				Surname:  "smith",
				Username: "fruitnvegAdmin",
				Password: []byte("123"),
			},
		},
	},
	"Reinhard Trucking": {
		{
			Client: client.Client{
				Name:              "Spar",
				AdminEmailAddress: "admin@spar.com",
			},
			AdminUser: user.User{
				Name:     "Sipho",
				Surname:  "Shezi",
				Username: "sparAdmin",
				Password: []byte("123"),
			},
		},
		{
			Client: client.Client{
				Name:              "Game",
				AdminEmailAddress: "admin@game.com",
			},
			AdminUser: user.User{
				Name:     "Michael",
				Surname:  "Black",
				Username: "gameAdmin",
				Password: []byte("123"),
			},
		},
	},
}

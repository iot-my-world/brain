package client

import (
	"gitlab.com/iotTracker/brain/party/client"
	"gitlab.com/iotTracker/brain/party/user"
)

type TestData struct {
	Client    client.Client
	AdminUser user.User
}

var EntitiesAndAdminUsersToCreate = map[string][]TestData{
	"Monteagle Logistics Limited": []TestData{
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
		},
	},
	"DHL": []TestData{
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
	"Reinhard Trucking": []TestData{
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

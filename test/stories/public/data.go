package public

import (
	"github.com/iot-my-world/brain/pkg/party/client"
	"github.com/iot-my-world/brain/pkg/party/company"
	humanUser "github.com/iot-my-world/brain/pkg/user/human"
	publicTestModule "github.com/iot-my-world/brain/test/modules/public"
)

var CompanyTestData = []publicTestModule.CompanyData{
	{
		Company: company.Company{
			Name:              "Public Sun Shop",
			AdminEmailAddress: "publicSunShopAdmin@sunShop.com",
		},
		AdminUser: humanUser.User{
			Name:     "Tom",
			Surname:  "Smith",
			Username: "publicSunShopAdmin",
			Password: []byte("123"),
			Roles:    make([]string, 0),
		},
		Users: []humanUser.User{
			{
				Name:         "sunShopUser1",
				Surname:      "test1",
				Username:     "sunShopUser1",
				Password:     []byte("123"),
				EmailAddress: "sunShopUser1@sunShop.com",
				Roles:        make([]string, 0),
			},
			{
				Name:         "sunShopUser2",
				Surname:      "test1",
				Username:     "sunShopUser2",
				Password:     []byte("123"),
				EmailAddress: "sunShopUser2@sunShop.com",
				Roles:        make([]string, 0),
			},
		},
	},
	{
		Company: company.Company{
			// Id:
			Name:              "Fastway",
			AdminEmailAddress: "fastwayAdmin@fastway.com",
			// ParentPartyType:
			// ParentId:
		},
		AdminUser: humanUser.User{
			Name:     "Simon",
			Surname:  "Rubin",
			Username: "fastwayAdmin",
			Password: []byte("123"),
			Roles:    make([]string, 0),
		},
		Users: []humanUser.User{
			{
				Name:         "fastwayUser1",
				Surname:      "test1",
				Username:     "fastwayUser1",
				Password:     []byte("123"),
				EmailAddress: "fastwayUser1@fastway.com",
				Roles:        make([]string, 0),
			},
			{
				Name:         "fastwayUser2",
				Surname:      "test2",
				Username:     "fastwayUser2",
				Password:     []byte("123"),
				EmailAddress: "fastwayUser2@fastway.com",
				Roles:        make([]string, 0),
			},
		},
	},
	{
		Company: company.Company{
			// Id:
			Name:              "Eskom",
			AdminEmailAddress: "eskomAdmin@eskom.com",
			// ParentPartyType:
			// ParentId:
		},
		AdminUser: humanUser.User{
			Name:     "Bob",
			Surname:  "Smith",
			Username: "eskomAdmin",
			Password: []byte("123"),
			Roles:    make([]string, 0),
		},
		Users: []humanUser.User{
			{
				Name:         "eskomUser1",
				Surname:      "test1",
				Username:     "eskomUser1",
				Password:     []byte("123"),
				EmailAddress: "eskomUser1@eskom.com",
				Roles:        make([]string, 0),
			},
			{
				Name:         "eskomUser2",
				Surname:      "test2",
				Username:     "eskomUser2",
				Password:     []byte("123"),
				EmailAddress: "eskomUser2@eskom.com",
				Roles:        make([]string, 0),
			},
		},
	},
}

var ClientTestData = []publicTestModule.ClientData{
	{
		Client: client.Client{
			Type:              "",
			Name:              "",
			AdminEmailAddress: "",
		},
		AdminUser: humanUser.User{
			Name:         "",
			Surname:      "",
			Username:     "",
			EmailAddress: "",
			Password:     []byte("123"),
			Roles:        make([]string, 0),
		},
		Users: []humanUser.User{
			{
				Name:         "",
				Surname:      "",
				Username:     "",
				EmailAddress: "",
				Password:     []byte("123"),
				Roles:        make([]string, 0),
			},
		},
	},
}

package public

import (
	"github.com/iot-my-world/brain/pkg/party/client"
	"github.com/iot-my-world/brain/pkg/party/company"
	humanUser "github.com/iot-my-world/brain/pkg/user/human"
	partyRegistrarAdministratorTestModule "github.com/iot-my-world/brain/test/modules/party/registrarAdministrator"
)

var CompanyTestData = []partyRegistrarAdministratorTestModule.CompanyData{
	{
		Company: company.Company{
			Name:              "Public Sun Shop",
			AdminEmailAddress: "publicSunShopAdmin@sunShop.com",
		},
		AdminUser: humanUser.User{
			Name:         "Tom",
			Surname:      "Smith",
			Username:     "publicSunShopAdmin",
			EmailAddress: "publicSunShopAdmin@sunShop.com",
			Password:     []byte("123"),
			Roles:        make([]string, 0),
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
			Name:         "Simon",
			Surname:      "Rubin",
			Username:     "fastwayAdmin",
			EmailAddress: "fastwayAdmin@fastway.com",
			Password:     []byte("123"),
			Roles:        make([]string, 0),
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
			Name:         "Bob",
			Surname:      "Smith",
			Username:     "eskomAdmin",
			EmailAddress: "eskomAdmin@eskom.com",
			Password:     []byte("123"),
			Roles:        make([]string, 0),
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

var ClientTestData = []partyRegistrarAdministratorTestModule.ClientData{
	{
		Client: client.Client{
			Type:              client.Individual,
			Name:              "Paul",
			AdminEmailAddress: "paul@gmail.com",
		},
		AdminUser: humanUser.User{
			Name:         "Paul",
			Surname:      "Smith",
			Username:     "paul",
			EmailAddress: "paul@gmail.com",
			Password:     []byte("123"),
			Roles:        make([]string, 0),
		},
		Users: []humanUser.User{
			{
				Name:         "Sandra",
				Surname:      "Smith",
				Username:     "sandra",
				EmailAddress: "sandra@gmail.com",
				Password:     []byte("123"),
				Roles:        make([]string, 0),
			},
		},
	},
	{
		Client: client.Client{
			Type:              client.Company,
			Name:              "Samson Logistics",
			AdminEmailAddress: "jacqui@gmail.com",
		},
		AdminUser: humanUser.User{
			Name:         "Jacqui",
			Surname:      "White",
			Username:     "jacqui",
			EmailAddress: "jacqui@gmail.com",
			Password:     []byte("123"),
			Roles:        make([]string, 0),
		},
		Users: []humanUser.User{
			{
				Name:         "Kyle",
				Surname:      "May",
				Username:     "kyleUser",
				EmailAddress: "kyle@gmail.com",
				Password:     []byte("123"),
				Roles:        make([]string, 0),
			},
		},
	},
}

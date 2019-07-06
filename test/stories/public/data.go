package public

import (
	"github.com/iot-my-world/brain/pkg/party/company"
	humanUser "github.com/iot-my-world/brain/pkg/user/human"
	companyTest "github.com/iot-my-world/brain/test/modules/party/company"
)

var TestData = []companyTest.Data{
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
			Name:              "DHL",
			AdminEmailAddress: "dhlAdmin@dhl.com",
			// ParentPartyType:
			// ParentId:
		},
		AdminUser: humanUser.User{
			Name:     "Simon",
			Surname:  "Rubin",
			Username: "dhlAdmin",
			Password: []byte("123"),
			Roles:    make([]string, 0),
		},
		Users: []humanUser.User{
			{
				Name:         "dhlUser1",
				Surname:      "test1",
				Username:     "dhlUser1",
				Password:     []byte("123"),
				EmailAddress: "dhlUser1@dhl.com",
				Roles:        make([]string, 0),
			},
			{
				Name:         "dhlUser2",
				Surname:      "test2",
				Username:     "dhlUser2",
				Password:     []byte("123"),
				EmailAddress: "dhlUser2@dhl.com",
				Roles:        make([]string, 0),
			},
		},
	},
	{
		Company: company.Company{
			// Id:
			Name:              "Reinhard Trucking",
			AdminEmailAddress: "reinhardAdmin@reinhard.com",
			// ParentPartyType:
			// ParentId:
		},
		AdminUser: humanUser.User{
			Name:     "Bob",
			Surname:  "Smith",
			Username: "reinhardAdmin",
			Password: []byte("123"),
			Roles:    make([]string, 0),
		},
		Users: []humanUser.User{
			{
				Name:         "reinhardUser1",
				Surname:      "test1",
				Username:     "reinhardUser1",
				Password:     []byte("123"),
				EmailAddress: "reinhardUser1@reinhard.com",
				Roles:        make([]string, 0),
			},
			{
				Name:         "reinhardUser2",
				Surname:      "test2",
				Username:     "reinhardUser2",
				Password:     []byte("123"),
				EmailAddress: "reinhardUser2@reinhard.com",
				Roles:        make([]string, 0),
			},
		},
	},
}

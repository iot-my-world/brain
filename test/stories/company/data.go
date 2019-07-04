package company

import (
	"github.com/iot-my-world/brain/pkg/party/company"
	companyTest "github.com/iot-my-world/brain/test/modules/party/company"
	humanUser "github.com/iot-my-world/brain/user/human"
)

var TestData = []companyTest.Data{
	{
		Company: company.Company{
			Name:              "Monteagle Logistics Limited",
			AdminEmailAddress: "monteagleAdmin@monteagle.com",
		},
		AdminUser: humanUser.User{
			Name:     "Murray",
			Surname:  "Griffin",
			Username: "monteagleAdmin",
			Password: []byte("123"),
			Roles:    make([]string, 0),
		},
		Users: []humanUser.User{
			{
				Name:         "monteagleUser1",
				Surname:      "test1",
				Username:     "monteagleUser1",
				Password:     []byte("123"),
				EmailAddress: "monteagleUser1@monteagle.com",
				Roles:        make([]string, 0),
			},
			{
				Name:         "monteagleUser2",
				Surname:      "test2",
				Username:     "monteagleUser2",
				Password:     []byte("123"),
				EmailAddress: "monteagleUser2@monteagle.com",
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

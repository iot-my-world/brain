package company

import (
	"gitlab.com/iotTracker/brain/party/company"
	humanUser "gitlab.com/iotTracker/brain/user/human"
)

type TestData struct {
	Company   company.Company
	AdminUser humanUser.User
	Users     []humanUser.User
}

var EntitiesAndAdminUsersToCreate = []TestData{
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
		},
		Users: []humanUser.User{
			{
				Name:         "monteagleUser1",
				Surname:      "test1",
				Username:     "monteagleUser1",
				Password:     []byte("123"),
				EmailAddress: "monteagleUser1@monteagle.com",
			},
			{
				Name:         "monteagleUser2",
				Surname:      "test2",
				Username:     "monteagleUser2",
				Password:     []byte("123"),
				EmailAddress: "monteagleUser2@monteagle.com",
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
		},
		Users: []humanUser.User{
			{
				Name:         "dhlUser1",
				Surname:      "test1",
				Username:     "dhlUser1",
				Password:     []byte("123"),
				EmailAddress: "dhlUser1@dhl.com",
			},
			{
				Name:         "dhlUser2",
				Surname:      "test2",
				Username:     "dhlUser2",
				Password:     []byte("123"),
				EmailAddress: "dhlUser2@dhl.com",
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
		},
		Users: []humanUser.User{
			{
				Name:         "reinhardUser1",
				Surname:      "test1",
				Username:     "reinhardUser1",
				Password:     []byte("123"),
				EmailAddress: "reinhardUser1@reinhard.com",
			},
			{
				Name:         "reinhardUser2",
				Surname:      "test2",
				Username:     "reinhardUser2",
				Password:     []byte("123"),
				EmailAddress: "reinhardUser2@reinhard.com",
			},
		},
	},
}

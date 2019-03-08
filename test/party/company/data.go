package company

import (
	"gitlab.com/iotTracker/brain/party/company"
	"gitlab.com/iotTracker/brain/party/user"
)

type TestData struct {
	Company   company.Company
	AdminUser user.User
	Users     []user.User
}

var EntitiesAndAdminUsersToCreate = []TestData{
	{
		Company: company.Company{
			Name:              "Monteagle Logistics Limited",
			AdminEmailAddress: "admin@monteagle.com",
		},
		AdminUser: user.User{
			Name:     "Murray",
			Surname:  "Griffin",
			Username: "monteagleAdmin",
			Password: []byte("123"),
		},
		Users: []user.User{
			{
				Name:         "monteagleUser1",
				Surname:      "test1",
				Username:     "monteagleTestUser1",
				Password:     []byte("123"),
				EmailAddress: "monteagleUser1@monteagle.com",
			},
			{
				Name:         "monteagleUser2",
				Surname:      "test2",
				Username:     "monteagleTestUser2",
				Password:     []byte("123"),
				EmailAddress: "monteagleUser2@monteagle.com",
			},
		},
	},
	{
		Company: company.Company{
			// Id:
			Name:              "DHL",
			AdminEmailAddress: "admin@dhl.com",
			// ParentPartyType:
			// ParentId:
		},
		AdminUser: user.User{
			Name:     "Simon",
			Surname:  "Rubin",
			Username: "dhlAdmin",
			Password: []byte("123"),
		},
		Users: []user.User{
			{
				Name:         "dhlUser1",
				Surname:      "test1",
				Username:     "dhlTestUser1",
				Password:     []byte("123"),
				EmailAddress: "dhlUser1@dhl.com",
			},
			{
				Name:         "dhlUser2",
				Surname:      "test2",
				Username:     "dhlTestUser2",
				Password:     []byte("123"),
				EmailAddress: "dhlUser2@dhl.com",
			},
		},
	},
	{
		Company: company.Company{
			// Id:
			Name:              "Reinhard Trucking",
			AdminEmailAddress: "admin@reinhard.com",
			// ParentPartyType:
			// ParentId:
		},
		AdminUser: user.User{
			Name:     "Bob",
			Surname:  "Smith",
			Username: "reinhardAdmin",
			Password: []byte("123"),
		},
		Users: []user.User{
			{
				Name:         "reinhardUser1",
				Surname:      "test1",
				Username:     "reinhardTestUser1",
				Password:     []byte("123"),
				EmailAddress: "reinhardUser1@reinhard.com",
			},
			{
				Name:         "reinhardUser2",
				Surname:      "test2",
				Username:     "reinhardTestUser2",
				Password:     []byte("123"),
				EmailAddress: "reinhardUser2@reinhard.com",
			},
		},
	},
}

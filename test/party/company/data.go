package company

import (
	"gitlab.com/iotTracker/brain/party/company"
	"gitlab.com/iotTracker/brain/party/user"
)

type TestData struct {
	Company   company.Company
	AdminUser user.User
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
	},
}

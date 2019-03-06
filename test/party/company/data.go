package company

import (
	"gitlab.com/iotTracker/brain/party/company"
	"gitlab.com/iotTracker/brain/party/user"
)

type EntityAdminUser struct {
	Company   company.Company
	AdminUser user.User
}

var EntitiesAndAdminUsersToCreate = []EntityAdminUser{
	{
		Company: company.Company{
			// Id:
			Name:              "Monteagle Logistics Limited",
			AdminEmailAddress: "monty@monteagle.com",
			// ParentPartyType:
			// ParentId:
		},
		AdminUser: user.User{
			Name:     "Murray",
			Surname:  "Griffin",
			Username: "murray",
			Password: []byte("123"),
		},
	},
	{
		Company: company.Company{
			// Id:
			Name:              "DHL",
			AdminEmailAddress: "dhlTest@dhl.com",
			// ParentPartyType:
			// ParentId:
		},
		AdminUser: user.User{
			Name:     "Simon",
			Surname:  "Rubin",
			Username: "simon",
			Password: []byte("123"),
		},
	},
	{
		Company: company.Company{
			// Id:
			Name:              "Reinhard Trucking",
			AdminEmailAddress: "reinhardTest@reinhard.com",
			// ParentPartyType:
			// ParentId:
		},
		AdminUser: user.User{
			Name:     "Bob",
			Surname:  "Smith",
			Username: "bob",
			Password: []byte("123"),
		},
	},
}

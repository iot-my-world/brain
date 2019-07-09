package company

import (
	"github.com/iot-my-world/brain/pkg/device/sigbug"
	"github.com/iot-my-world/brain/pkg/party"
	"github.com/iot-my-world/brain/pkg/party/company"
	"github.com/iot-my-world/brain/pkg/search/identifier/id"
	humanUser "github.com/iot-my-world/brain/pkg/user/human"
	companyTest "github.com/iot-my-world/brain/test/modules/party/company"
)

var TestData = []struct {
	CompanyTestData companyTest.Data
	SigbugDevices   []sigbug.Sigbug
}{
	{
		CompanyTestData: companyTest.Data{
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
		SigbugDevices: []sigbug.Sigbug{
			// device owned by company, not yet assigned to client
			{
				DeviceId:       "Monteagle Logistics Limited",
				OwnerPartyType: party.Company,
				OwnerId: id.Identifier{
					Id: "Monteagle Logistics Limited", // populated correctly on creation
				},
				AssignedPartyType: "",
				AssignedId:        id.Identifier{},
			},
			// device owned by company assigned to a client
			{
				DeviceId:       "Monteagle Logistics Limited",
				OwnerPartyType: party.Company,
				OwnerId: id.Identifier{
					Id: "Monteagle Logistics Limited", // populated correctly on creation
				},
				AssignedPartyType: party.Client,
				AssignedId: id.Identifier{
					Id: "Pick 'n Pay", // populated correctly on creation
				},
			},
		},
	},

	{
		CompanyTestData: companyTest.Data{
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
		SigbugDevices: []sigbug.Sigbug{
			// device owned by company, not yet assigned to client
			{
				DeviceId:       "DHL",
				OwnerPartyType: party.Company,
				OwnerId: id.Identifier{
					Id: "DHL", // populated correctly on creation
				},
				AssignedPartyType: "",
				AssignedId:        id.Identifier{},
			},
			// device owned by company assigned to a client
			{
				DeviceId:       "DHL",
				OwnerPartyType: party.Company,
				OwnerId: id.Identifier{
					Id: "DHL", // populated correctly on creation
				},
				AssignedPartyType: party.Client,
				AssignedId: id.Identifier{
					Id: "Fruit 'n Veg", // populated correctly on creation
				},
			},
		},
	},

	{
		CompanyTestData: companyTest.Data{
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
		SigbugDevices: []sigbug.Sigbug{
			// device owned by company, not yet assigned to client
			{
				DeviceId:       "Reinhard Trucking",
				OwnerPartyType: party.Company,
				OwnerId: id.Identifier{
					Id: "Reinhard Trucking", // populated correctly on creation
				},
				AssignedPartyType: "",
				AssignedId:        id.Identifier{},
			},
			// device owned by company assigned to a client
			{
				DeviceId:       "Reinhard Trucking",
				OwnerPartyType: party.Company,
				OwnerId: id.Identifier{
					Id: "Reinhard Trucking", // populated correctly on creation
				},
				AssignedPartyType: party.Client,
				AssignedId: id.Identifier{
					Id: "Spar", // populated correctly on creation
				},
			},
		},
	},
}

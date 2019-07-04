package company

import (
	"github.com/iot-my-world/brain/pkg/party"
	"github.com/iot-my-world/brain/pkg/search/identifier/id"
)

type Company struct {
	Id   string `json:"id" bson:"id"`
	Name string `json:"name" bson:"name"`
	// The email address which will be used to invite the admin
	// user of the company
	// I.e. the first user of the system from the company
	AdminEmailAddress string `json:"adminEmailAddress" bson:"adminEmailAddress"`

	ParentPartyType party.Type    `json:"parentPartyType" bson:"parentPartyType"`
	ParentId        id.Identifier `json:"parentId" bson:"parentId"`
}

type wrapped struct {
	Id   string `json:"id" bson:"id"`
	Name string `json:"name" bson:"name"`
	// The email address which will be used to invite the admin
	// user of the company
	// I.e. the first user of the system from the company
	AdminEmailAddress string `json:"adminEmailAddress" bson:"adminEmailAddress"`

	ParentPartyType party.Type    `json:"parentPartyType" bson:"parentPartyType"`
	ParentId        id.Identifier `json:"parentId" bson:"parentId"`
}

func (c Company) Details() party.Details {
	return party.Details{
		ParentDetail: party.ParentDetail{
			ParentId:        c.ParentId,
			ParentPartyType: c.ParentPartyType,
		},
		Detail: party.Detail{
			PartyId:   id.Identifier{Id: c.Id},
			PartyType: party.Company,
		},
	}
}

func (c *Company) SetId(id string) {
	c.Id = id
}

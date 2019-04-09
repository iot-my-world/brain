package company

import (
	"gitlab.com/iotTracker/brain/party"
	"gitlab.com/iotTracker/brain/search/identifier"
	"gitlab.com/iotTracker/brain/search/identifier/id"
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

func (c Company) SetId(id string) {
	c.Id = id
}

func (c Company) ValidIdentifier(id identifier.Identifier) bool {
	if id == nil {
		return false
	}

	switch id.Type() {
	case identifier.Id, identifier.EmailAddress, identifier.AdminEmailAddress:
		return true
	default:
		return false
	}
}

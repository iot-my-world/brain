package company

import (
	"github.com/satori/go.uuid"
	"gitlab.com/iotTracker/brain/party"
	"gitlab.com/iotTracker/brain/search/identifier/id"
)

// Company is the model for the company entities in the system
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

// Details returns the party details of the company party
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

func (c Company) SetId() error {
	newId, err := uuid.NewV4()
	c.Id = newId.String()
	if err != nil {
		return err
	}
	return nil
}
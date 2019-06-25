package individual

import (
	"github.com/iot-my-world/brain/party"
	"github.com/iot-my-world/brain/search/identifier/id"
)

// Individual is the model for the client entities in the system
type Individual struct {
	Id string `json:"id" bson:"id"`

	Name string `json:"name" bson:"name"`

	// The email address which will be used to invite the admin
	// user of the client
	// I.e. the first user of the system from the individual
	AdminEmailAddress string `json:"adminEmailAddress" bson:"adminEmailAddress"`

	ParentPartyType party.Type    `json:"parentPartyType" bson:"parentPartyType"`
	ParentId        id.Identifier `json:"parentId" bson:"parentId"`
}

// Details returns the party details of the individual party
func (i Individual) Details() party.Details {
	return party.Details{
		ParentDetail: party.ParentDetail{
			ParentId:        i.ParentId,
			ParentPartyType: i.ParentPartyType,
		},
		Detail: party.Detail{
			PartyId:   id.Identifier{Id: i.Id},
			PartyType: party.Individual,
		},
	}
}

func (i *Individual) SetId(id string) {
	i.Id = id
}

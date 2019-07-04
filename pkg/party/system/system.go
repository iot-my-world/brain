package system

import (
	"github.com/iot-my-world/brain/pkg/party"
	"github.com/iot-my-world/brain/pkg/search/identifier/id"
)

// System is the model for the system entities in the system
type System struct {
	Id                string `json:"id" bson:"id"`
	Name              string `json:"name" bson:"name"`
	AdminEmailAddress string `json:"adminEmailAddress" bson:"adminEmailAddress"`
}

// Details returns the party details of the system party
func (s System) Details() party.Details {
	return party.Details{
		ParentDetail: party.ParentDetail{
			ParentId:        id.Identifier{Id: s.Id},
			ParentPartyType: party.System,
		},
		Detail: party.Detail{
			PartyId:   id.Identifier{Id: s.Id},
			PartyType: party.System,
		},
	}
}

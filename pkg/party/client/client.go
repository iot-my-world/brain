package client

import (
	"github.com/iot-my-world/brain/pkg/party"
	"github.com/iot-my-world/brain/pkg/search/identifier/id"
)

type Type string

const Individual Type = "Individual"
const Company Type = "Company"

// Client is the model for the client entities in the system
type Client struct {
	Id string `json:"id" bson:"id"`

	// Client Type, Individual or Company
	Type Type `json:"type" bson:"type"`

	Name string `json:"name" bson:"name"`

	// The email address which will be used to invite the admin
	// user of the client
	// I.e. the first user of the system from the client
	AdminEmailAddress string `json:"adminEmailAddress" bson:"adminEmailAddress"`

	ParentPartyType party.Type    `json:"parentPartyType" bson:"parentPartyType"`
	ParentId        id.Identifier `json:"parentId" bson:"parentId"`
}

// Details returns the party details of the client party
func (c Client) Details() party.Details {
	return party.Details{
		ParentDetail: party.ParentDetail{
			ParentId:        c.ParentId,
			ParentPartyType: c.ParentPartyType,
		},
		Detail: party.Detail{
			PartyId:   id.Identifier{Id: c.Id},
			PartyType: party.Client,
		},
	}
}

func (c *Client) SetId(id string) {
	c.Id = id
}

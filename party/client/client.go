package client

import (
	"gitlab.com/iotTracker/brain/party"
	"gitlab.com/iotTracker/brain/search/identifier/id"
)

type Client struct {
	Id              string        `json:"id" bson:"id"`

	Name            string        `json:"name" bson:"name"`

	// The email address which will be used to invite the admin
	// user of the company
	// I.e. the first user of the system from the company
	AdminEmailAddress string `json:"adminEmailAddress" bson:"adminEmailAddress"`

	ParentPartyType party.Type    `json:"parentPartyType" bson:"parentPartyType"`
	ParentId        id.Identifier `json:"parentId" bson:"parentId"`
}

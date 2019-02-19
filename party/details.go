package party

import "gitlab.com/iotTracker/brain/search/identifier/id"

type Details struct {
	PartyType Type          `json:"partyType"`
	PartyId   id.Identifier `json:"partyId"`
}

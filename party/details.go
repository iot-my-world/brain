package party

import "gitlab.com/iotTracker/brain/search/identifier/id"

type Details struct {
	ParentPartyType Type          `json:"parentPartyType"`
	ParentId        id.Identifier `json:"parentId"`
	PartyType       Type          `json:"partyType"`
	PartyId         id.Identifier `json:"partyId"`
}

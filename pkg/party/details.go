package party

import "github.com/iot-my-world/brain/search/identifier/id"

type Details struct {
	Detail
	ParentDetail
}

type Detail struct {
	PartyType Type          `json:"partyType"`
	PartyId   id.Identifier `json:"partyId"`
}

type ParentDetail struct {
	ParentPartyType Type          `json:"parentPartyType"`
	ParentId        id.Identifier `json:"parentId"`
}

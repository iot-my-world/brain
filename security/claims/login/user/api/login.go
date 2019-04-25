package api

import (
	"gitlab.com/iotTracker/brain/party"
	"gitlab.com/iotTracker/brain/search/identifier/id"
	"gitlab.com/iotTracker/brain/security/claims"
	"time"
)

type Login struct {
	UserId          id.Identifier `json:"userId"`
	IssueTime       int64         `json:"issueTime"`
	ExpirationTime  int64         `json:"expirationTime"`
	ParentPartyType party.Type    `json:"parentPartyType"`
	ParentId        id.Identifier `json:"parentId"`
	PartyType       party.Type    `json:"partyType"`
	PartyId         id.Identifier `json:"partyId"`
}

func (l Login) Type() claims.Type {
	return claims.APIUserLogin
}

func (l Login) Expired() bool {
	return time.Now().UTC().After(time.Unix(l.ExpirationTime, 0).UTC())
}

func (l Login) PartyDetails() party.Details {
	return party.Details{
		Detail: party.Detail{
			PartyType: l.PartyType,
			PartyId:   l.PartyId,
		},
		ParentDetail: party.ParentDetail{
			ParentPartyType: l.ParentPartyType,
			ParentId:        l.ParentId,
		},
	}
}

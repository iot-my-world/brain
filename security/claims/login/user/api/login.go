package api

import (
	"github.com/iot-my-world/brain/pkg/party"
	"github.com/iot-my-world/brain/search/identifier/id"
	"github.com/iot-my-world/brain/security/claims"
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

func (l Login) TimeToExpiry() time.Duration {
	return time.Unix(l.ExpirationTime, 0).UTC().Sub(time.Now().UTC())
}

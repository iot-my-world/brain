package claims

import (
	"gitlab.com/iotTracker/brain/search/identifier/id"
	"gitlab.com/iotTracker/brain/party"
)

type LoginClaims struct {
	UserId         id.Identifier `json:"userId"`
	IssuedAtTime   int64         `json:"issueTime"`
	ExpirationTime int64         `json:"expiry"`
	PartyType      party.Type    `json:"partyType"`
	PartyId        id.Identifier `json:"partyId"`
}

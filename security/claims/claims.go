package claims

import (
	"gitlab.com/iotTracker/brain/search/identifier/id"
	"gitlab.com/iotTracker/brain/party"
	"time"
)

const ValidTime = 90*time.Minute

type Claims struct {
	UserId         id.Identifier `json:"userId"`
	IssueTime      int64         `json:"issueTime"`
	ExpirationTime int64         `json:"expirationTime"`
	PartyType      party.Type    `json:"partyType"`
	PartyId        id.Identifier `json:"partyId"`
}

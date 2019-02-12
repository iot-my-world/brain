package claims

import (
	"gitlab.com/iotTracker/brain/party"
	"gitlab.com/iotTracker/brain/search/identifier/id"
	"time"
)

const ValidTime = 90 * time.Minute

type LoginClaims struct {
	UserId         id.Identifier `json:"userId"`
	IssueTime      int64         `json:"issueTime"`
	ExpirationTime int64         `json:"expirationTime"`
	PartyType      party.Type    `json:"partyType"`
	PartyId        id.Identifier `json:"partyId"`
}

type RegisterCompanyAdminUserClaims struct {
	IssueTime      int64         `json:"issueTime"`
	ExpirationTime int64         `json:"expirationTime"`
	PartyId        id.Identifier `json:"partyId"`
}

type RegisterClientAdminUserClaims struct {
	IssueTime      int64         `json:"issueTime"`
	ExpirationTime int64         `json:"expirationTime"`
	PartyId        id.Identifier `json:"partyId"`
}

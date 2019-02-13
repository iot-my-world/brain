package registerCompanyAdminUser

import (
	"gitlab.com/iotTracker/brain/party"
	"gitlab.com/iotTracker/brain/search/identifier/id"
	"gitlab.com/iotTracker/brain/security/claims"
	"time"
)

type RegisterCompanyAdminUser struct {
	IssueTime      int64         `json:"issueTime"`
	ExpirationTime int64         `json:"expirationTime"`
	PartyType      party.Type    `json:"partyType"`
	PartyId        id.Identifier `json:"partyId"`
}

func (r RegisterCompanyAdminUser) Type() claims.Type {
	return claims.RegisterCompanyAdminUser
}

func (r RegisterCompanyAdminUser) Expired() bool {
	return time.Now().UTC().After(time.Unix(r.ExpirationTime, 0).UTC())
}

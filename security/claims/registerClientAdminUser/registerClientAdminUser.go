package registerClientAdminUser

import (
	"gitlab.com/iotTracker/brain/party"
	"gitlab.com/iotTracker/brain/search/identifier/id"
	"gitlab.com/iotTracker/brain/security/claims"
	"gitlab.com/iotTracker/brain/security/permission/api"
	"time"
)

type RegisterClientAdminUser struct {
	IssueTime      int64         `json:"issueTime"`
	ExpirationTime int64         `json:"expirationTime"`
	PartyType      party.Type    `json:"partyType"`
	PartyId        id.Identifier `json:"partyId"`
	EmailAddress   string        `json:"emailAddress"`
}

func (r RegisterClientAdminUser) Type() claims.Type {
	return claims.RegisterClientAdminUser
}

func (r RegisterClientAdminUser) Expired() bool {
	return time.Now().UTC().After(time.Unix(r.ExpirationTime, 0).UTC())
}

func (r RegisterClientAdminUser) PartyDetails() party.Details {
	return party.Details{
		PartyType: r.PartyType,
		PartyId:   r.PartyId,
	}
}

// permissions granted by having a valid set of these claims
var GrantedAPIPermissions = []api.Permission{
	api.UserRecordHandlerValidate,             // Ability to validate users
	api.PartyRegistrarRegisterClientAdminUser, // Ability to register self
}

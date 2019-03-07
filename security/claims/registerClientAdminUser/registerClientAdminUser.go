package registerClientAdminUser

import (
	"gitlab.com/iotTracker/brain/party"
	"gitlab.com/iotTracker/brain/search/identifier/id"
	"gitlab.com/iotTracker/brain/security/claims"
	"gitlab.com/iotTracker/brain/security/permission/api"
	"time"
	"gitlab.com/iotTracker/brain/party/user"
)

type RegisterClientAdminUser struct {
	IssueTime       int64         `json:"issueTime"`
	ExpirationTime  int64         `json:"expirationTime"`
	EmailAddress    string        `json:"emailAddress"`
	ParentPartyType party.Type    `json:"parentPartyType"`
	ParentId        id.Identifier `json:"parentId"`
	PartyType       party.Type    `json:"partyType"`
	PartyId         id.Identifier `json:"partyId"`
	User            user.User     `json:"user"`
}

func (r RegisterClientAdminUser) Type() claims.Type {
	return claims.RegisterClientAdminUser
}

func (r RegisterClientAdminUser) Expired() bool {
	return time.Now().UTC().After(time.Unix(r.ExpirationTime, 0).UTC())
}

func (r RegisterClientAdminUser) PartyDetails() party.Details {
	return party.Details{
		ParentPartyType: r.ParentPartyType,
		ParentId:        r.ParentId,
		PartyType:       r.PartyType,
		PartyId:         r.PartyId,
	}
}

// permissions granted by having a valid set of these claims
var GrantedAPIPermissions = []api.Permission{
	api.UserRecordHandlerValidate,             // Ability to validate users
	api.PartyRegistrarRegisterClientAdminUser, // Ability to register self
}

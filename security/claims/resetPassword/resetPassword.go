package resetPassword

import (
	"gitlab.com/iotTracker/brain/party"
	"gitlab.com/iotTracker/brain/search/identifier/id"
	"gitlab.com/iotTracker/brain/security/claims"
	"gitlab.com/iotTracker/brain/security/permission/api"
	"time"
)

type ResetPassword struct {
	UserId          id.Identifier `json:"userId"`
	IssueTime       int64         `json:"issueTime"`
	ExpirationTime  int64         `json:"expirationTime"`
	ParentPartyType party.Type    `json:"parentPartyType"`
	ParentId        id.Identifier `json:"parentId"`
	PartyType       party.Type    `json:"partyType"`
	PartyId         id.Identifier `json:"partyId"`
}

func (r ResetPassword) Type() claims.Type {
	return claims.ResetPassword
}

func (r ResetPassword) Expired() bool {
	return time.Now().UTC().After(time.Unix(r.ExpirationTime, 0).UTC())
}

func (r ResetPassword) PartyDetails() party.Details {
	return party.Details{
		Detail: party.Detail{
			PartyType: r.PartyType,
			PartyId:   r.PartyId,
		},
		ParentDetail: party.ParentDetail{
			ParentPartyType: r.ParentPartyType,
			ParentId:        r.ParentId,
		},
	}
}

// permissions granted by having a valid set of these claims
var GrantedAPIPermissions = []api.Permission{
	api.UserAdministratorSetPassword,
}

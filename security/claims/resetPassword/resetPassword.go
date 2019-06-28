package resetPassword

import (
	"github.com/iot-my-world/brain/party"
	"github.com/iot-my-world/brain/search/identifier/id"
	"github.com/iot-my-world/brain/security/claims"
	"github.com/iot-my-world/brain/security/permission/api"
	humanUserAdministrator "github.com/iot-my-world/brain/user/human/administrator"
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

func (r ResetPassword) TimeToExpiry() time.Duration {
	return time.Unix(r.ExpirationTime, 0).UTC().Sub(time.Now().UTC())
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
	humanUserAdministrator.SetPasswordService,
}

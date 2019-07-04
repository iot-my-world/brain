package resetPassword

import (
	"github.com/iot-my-world/brain/pkg/party"
	"github.com/iot-my-world/brain/pkg/search/identifier/id"
	claims2 "github.com/iot-my-world/brain/pkg/security/claims"
	api2 "github.com/iot-my-world/brain/pkg/security/permission/api"
	humanUserAdministrator "github.com/iot-my-world/brain/pkg/user/human/administrator"
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

func (r ResetPassword) Type() claims2.Type {
	return claims2.ResetPassword
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
var GrantedAPIPermissions = []api2.Permission{
	humanUserAdministrator.SetPasswordService,
}

package registerClientAdminUser

import (
	"github.com/iot-my-world/brain/party"
	partyRegistrar "github.com/iot-my-world/brain/party/registrar"
	"github.com/iot-my-world/brain/search/identifier/id"
	"github.com/iot-my-world/brain/security/claims"
	"github.com/iot-my-world/brain/security/permission/api"
	humanUser "github.com/iot-my-world/brain/user/human"
	humanUserValidator "github.com/iot-my-world/brain/user/human/validator"
	"time"
)

type RegisterClientAdminUser struct {
	IssueTime       int64          `json:"issueTime"`
	ExpirationTime  int64          `json:"expirationTime"`
	ParentPartyType party.Type     `json:"parentPartyType"`
	ParentId        id.Identifier  `json:"parentId"`
	PartyType       party.Type     `json:"partyType"`
	PartyId         id.Identifier  `json:"partyId"`
	User            humanUser.User `json:"user"`
}

func (r RegisterClientAdminUser) Type() claims.Type {
	return claims.RegisterClientAdminUser
}

func (r RegisterClientAdminUser) Expired() bool {
	return time.Now().UTC().After(time.Unix(r.ExpirationTime, 0).UTC())
}

func (r RegisterClientAdminUser) TimeToExpiry() time.Duration {
	return time.Unix(r.ExpirationTime, 0).UTC().Sub(time.Now().UTC())
}

func (r RegisterClientAdminUser) PartyDetails() party.Details {
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
	humanUserValidator.ValidateService,            // Ability to validate users
	partyRegistrar.RegisterClientAdminUserService, // Ability to register self
}

package registerCompanyAdminUser

import (
	"github.com/iot-my-world/brain/pkg/party"
	partyRegistrar "github.com/iot-my-world/brain/pkg/party/registrar"
	"github.com/iot-my-world/brain/pkg/search/identifier/id"
	claims "github.com/iot-my-world/brain/pkg/security/claims"
	apiPermission "github.com/iot-my-world/brain/pkg/security/permission/api"
	humanUser "github.com/iot-my-world/brain/pkg/user/human"
	humanUserValidator "github.com/iot-my-world/brain/pkg/user/human/validator"
	"time"
)

type RegisterCompanyAdminUser struct {
	IssueTime       int64          `json:"issueTime"`
	ExpirationTime  int64          `json:"expirationTime"`
	User            humanUser.User `json:"user"`
	ParentPartyType party.Type     `json:"parentPartyType"`
	ParentId        id.Identifier  `json:"parentId"`
	PartyType       party.Type     `json:"partyType"`
	PartyId         id.Identifier  `json:"partyId"`
}

func (r RegisterCompanyAdminUser) Type() claims.Type {
	return claims.RegisterCompanyAdminUser
}

func (r RegisterCompanyAdminUser) Expired() bool {
	return time.Now().UTC().After(time.Unix(r.ExpirationTime, 0).UTC())
}

func (r RegisterCompanyAdminUser) TimeToExpiry() time.Duration {
	return time.Unix(r.ExpirationTime, 0).UTC().Sub(time.Now().UTC())
}

func (r RegisterCompanyAdminUser) PartyDetails() party.Details {
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
var GrantedAPIPermissions = []apiPermission.Permission{
	humanUserValidator.ValidateService,             // Ability to validate users
	partyRegistrar.RegisterCompanyAdminUserService, // Ability to register self
}

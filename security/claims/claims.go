package claims

import (
	"gitlab.com/iotTracker/brain/party"
	"gitlab.com/iotTracker/brain/search/identifier/id"
	"time"
)

const ValidTime = 90 * time.Minute

type Type string

const Login Type = "Login"
const RegisterCompanyAdminUser Type = "RegisterCompanyAdminUser"
const RegisterClientAdminUser Type = "RegisterClientAdminUser"

type LoginClaims struct {
	Type           Type          `json:"type"`
	UserId         id.Identifier `json:"userId"`
	IssueTime      int64         `json:"issueTime"`
	ExpirationTime int64         `json:"expirationTime"`
	PartyType      party.Type    `json:"partyType"`
	PartyId        id.Identifier `json:"partyId"`
}

type RegisterCompanyAdminUserClaims struct {
	Type           Type          `json:"type"`
	IssueTime      int64         `json:"issueTime"`
	ExpirationTime int64         `json:"expirationTime"`
	PartyType      party.Type    `json:"partyType"`
	PartyId        id.Identifier `json:"partyId"`
}

type RegisterClientAdminUserClaims struct {
	Type           Type          `json:"type"`
	IssueTime      int64         `json:"issueTime"`
	ExpirationTime int64         `json:"expirationTime"`
	PartyType      party.Type    `json:"partyType"`
	PartyId        id.Identifier `json:"partyId"`
}

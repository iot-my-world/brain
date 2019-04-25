package claims

import (
	"gitlab.com/iotTracker/brain/party"
)

type Type string

const HumanUserLogin Type = "HumanUserLogin"
const APIUserLogin Type = "APIUserLogin"
const RegisterCompanyAdminUser Type = "RegisterCompanyAdminUser"
const RegisterCompanyUser Type = "RegisterCompanyUser"
const RegisterClientAdminUser Type = "RegisterClientAdminUser"
const RegisterClientUser Type = "RegisterClientUser"
const ResetPassword Type = "ResetPassword"

type Claims interface {
	Type() Type
	Expired() bool
	PartyDetails() party.Details
}

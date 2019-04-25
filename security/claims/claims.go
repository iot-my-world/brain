package claims

import (
	"gitlab.com/iotTracker/brain/party"
	"time"
)

type Type string

const HumanUserLogin Type = "HumanUserLogin"
const RegisterCompanyAdminUser Type = "RegisterCompanyAdminUser"
const RegisterCompanyUser Type = "RegisterCompanyUser"
const RegisterClientAdminUser Type = "RegisterClientAdminUser"
const RegisterClientUser Type = "RegisterClientUser"
const ResetPassword Type = "ResetPassword"

const ValidTime = 90 * time.Minute

type Claims interface {
	Type() Type
	Expired() bool
	PartyDetails() party.Details
}

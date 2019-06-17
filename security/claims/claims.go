package claims

import (
	"github.com/iot-my-world/brain/party"
	"time"
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
	TimeToExpiry() time.Duration
	PartyDetails() party.Details
}

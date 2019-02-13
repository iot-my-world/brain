package claims

import (
	"time"
)

type Type string

const Login Type = "Login"
const RegisterCompanyAdminUser Type = "RegisterCompanyAdminUser"
const RegisterClientAdminUser Type = "RegisterClientAdminUser"

const ValidTime = 90 * time.Minute

type Claims interface {
	Type() Type
	Expired() bool
}

package validator

import (
	"github.com/iot-my-world/brain/action"
	sf0012 "github.com/iot-my-world/brain/pkg/tracker/sf001"
	"github.com/iot-my-world/brain/security/claims"
	"github.com/iot-my-world/brain/security/permission/api"
	"github.com/iot-my-world/brain/validate/reasonInvalid"
)

type Validator interface {
	Validate(request *ValidateRequest) (*ValidateResponse, error)
}

const ServiceProvider = "SF001Tracker-Validator"
const ValidateService = ServiceProvider + ".Validate"

var SystemUserPermissions = []api.Permission{
	ValidateService,
}

var CompanyAdminUserPermissions = make([]api.Permission, 0)

var CompanyUserPermissions = make([]api.Permission, 0)

var ClientAdminUserPermissions = make([]api.Permission, 0)

var ClientUserPermissions = make([]api.Permission, 0)

type ValidateRequest struct {
	Claims claims.Claims
	SF001  sf0012.SF001
	Action action.Action
}

type ValidateResponse struct {
	ReasonsInvalid []reasonInvalid.ReasonInvalid
}

package validator

import (
	"github.com/iot-my-world/brain/action"
	"github.com/iot-my-world/brain/pkg/user/human"
	"github.com/iot-my-world/brain/security/claims"
	"github.com/iot-my-world/brain/security/permission/api"
	"github.com/iot-my-world/brain/validate/reasonInvalid"
)

type Validator interface {
	Validate(request *ValidateRequest) (*ValidateResponse, error)
}

const ServiceProvider = "HumanUser-Validator"
const ValidateService = ServiceProvider + ".Validate"

var SystemUserPermissions = make([]api.Permission, 0)

var CompanyAdminUserPermissions = []api.Permission{
	ValidateService,
}

var CompanyUserPermissions = make([]api.Permission, 0)

var ClientAdminUserPermissions = []api.Permission{
	ValidateService,
}

var ClientUserPermissions = make([]api.Permission, 0)

type ValidateRequest struct {
	Claims claims.Claims
	User   human.User
	Action action.Action
}

type ValidateResponse struct {
	ReasonsInvalid []reasonInvalid.ReasonInvalid
}

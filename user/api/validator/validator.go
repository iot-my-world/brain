package validator

import (
	"github.com/iot-my-world/brain/action"
	"github.com/iot-my-world/brain/security/claims"
	"github.com/iot-my-world/brain/security/permission/api"
	apiUser "github.com/iot-my-world/brain/user/api"
	"github.com/iot-my-world/brain/validate/reasonInvalid"
)

type Validator interface {
	Validate(request *ValidateRequest) (*ValidateResponse, error)
}

const ServiceProvider = "APIUser-Validator"
const ValidateService = ServiceProvider + ".Validate"

var SystemUserPermissions = []api.Permission{
	ValidateService,
}

var CompanyAdminUserPermissions = make([]api.Permission, 0)

var CompanyUserPermissions = make([]api.Permission, 0)

var ClientAdminUserPermissions = make([]api.Permission, 0)

type ValidateRequest struct {
	Claims claims.Claims
	User   apiUser.User
	Action action.Action
}

type ValidateResponse struct {
	ReasonsInvalid []reasonInvalid.ReasonInvalid
}

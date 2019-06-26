package validator

import (
	"github.com/iot-my-world/brain/action"
	"github.com/iot-my-world/brain/party/individual"
	"github.com/iot-my-world/brain/security/claims"
	"github.com/iot-my-world/brain/security/permission/api"
	"github.com/iot-my-world/brain/validate/reasonInvalid"
)

type Validator interface {
	Validate(request *ValidateRequest) (*ValidateResponse, error)
}

const ServiceProvider = "Individual-Validator"
const ValidateService = ServiceProvider + ".Validate"

var SystemUserPermissions = make([]api.Permission, 0)

var CompanyAdminUserPermissions = make([]api.Permission, 0)

var CompanyUserPermissions = make([]api.Permission, 0)

var ClientAdminUserPermissions = make([]api.Permission, 0)

var ClientUserPermissions = make([]api.Permission, 0)

type ValidateRequest struct {
	Claims     claims.Claims
	Individual individual.Individual
	Action     action.Action
}

type ValidateResponse struct {
	ReasonsInvalid []reasonInvalid.ReasonInvalid
}

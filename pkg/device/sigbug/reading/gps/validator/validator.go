package validator

import (
	"github.com/iot-my-world/brain/pkg/action"
	sigbugGPSReading "github.com/iot-my-world/brain/pkg/device/sigbug/reading/gps"
	"github.com/iot-my-world/brain/pkg/security/claims"
	"github.com/iot-my-world/brain/pkg/security/permission/api"
	"github.com/iot-my-world/brain/pkg/validate/reasonInvalid"
)

type Validator interface {
	Validate(request *ValidateRequest) (*ValidateResponse, error)
}

const ServiceProvider = "SigbugGPSReading-Validator"
const ValidateService = ServiceProvider + ".Validate"

var SystemUserPermissions = []api.Permission{
	ValidateService,
}

var CompanyAdminUserPermissions = make([]api.Permission, 0)

var CompanyUserPermissions = make([]api.Permission, 0)

var ClientAdminUserPermissions = make([]api.Permission, 0)

var ClientUserPermissions = make([]api.Permission, 0)

type ValidateRequest struct {
	Claims  claims.Claims
	Reading sigbugGPSReading.Reading
	Action  action.Action
}

type ValidateResponse struct {
	ReasonsInvalid []reasonInvalid.ReasonInvalid
}

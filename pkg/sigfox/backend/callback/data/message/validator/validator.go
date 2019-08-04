package validator

import (
	"github.com/iot-my-world/brain/pkg/action"
	"github.com/iot-my-world/brain/pkg/security/claims"
	"github.com/iot-my-world/brain/pkg/security/permission/api"
	sigfoxBackendDataCallbackMessage "github.com/iot-my-world/brain/pkg/sigfox/backend/callback/data/message"
	"github.com/iot-my-world/brain/pkg/validate/reasonInvalid"
)

type Validator interface {
	Validate(request *ValidateRequest) (*ValidateResponse, error)
}

const ServiceProvider = "MessageDevice-Validator"
const ValidateService = ServiceProvider + ".Validate"

var SystemUserPermissions = []api.Permission{
	ValidateService,
}

var CompanyAdminUserPermissions = []api.Permission{
	ValidateService,
}

var CompanyUserPermissions = make([]api.Permission, 0)

var ClientAdminUserPermissions = []api.Permission{
	ValidateService,
}

var ClientUserPermissions = make([]api.Permission, 0)

type ValidateRequest struct {
	Claims  claims.Claims
	Message sigfoxBackendDataCallbackMessage.Message
	Action  action.Action
}

type ValidateResponse struct {
	ReasonsInvalid []reasonInvalid.ReasonInvalid
}

package validator

import (
	"github.com/iot-my-world/brain/pkg/action"
	client2 "github.com/iot-my-world/brain/pkg/party/client"
	"github.com/iot-my-world/brain/pkg/validate/reasonInvalid"
	"github.com/iot-my-world/brain/security/claims"
	"github.com/iot-my-world/brain/security/permission/api"
)

type Validator interface {
	Validate(request *ValidateRequest) (*ValidateResponse, error)
}

const ServiceProvider = "Client-Validator"
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
	Client client2.Client
	Action action.Action
}

type ValidateResponse struct {
	ReasonsInvalid []reasonInvalid.ReasonInvalid
}

package administrator

import (
	"github.com/iot-my-world/brain/party/individual"
	"github.com/iot-my-world/brain/security/claims"
	"github.com/iot-my-world/brain/security/permission/api"
)

type Administrator interface {
	Create(request *CreateRequest) (*CreateResponse, error)
	UpdateAllowedFields(request *UpdateAllowedFieldsRequest) (*UpdateAllowedFieldsResponse, error)
}

const ServiceProvider = "Individual-Administrator"
const UpdateAllowedFieldsService = ServiceProvider + ".UpdateAllowedFields"
const CreateService = ServiceProvider + ".Create"

var SystemUserPermissions = make([]api.Permission, 0)

var CompanyAdminUserPermissions = make([]api.Permission, 0)

var CompanyUserPermissions = make([]api.Permission, 0)

var ClientAdminUserPermissions = make([]api.Permission, 0)

var ClientUserPermissions = make([]api.Permission, 0)

type CreateRequest struct {
	Claims     claims.Claims
	Individual individual.Individual
}

type CreateResponse struct {
	Individual individual.Individual
}

type UpdateAllowedFieldsRequest struct {
	Claims     claims.Claims
	Individual individual.Individual
}

type UpdateAllowedFieldsResponse struct {
	Individual individual.Individual
}

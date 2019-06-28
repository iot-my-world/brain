package recordHandler

import (
	"github.com/iot-my-world/brain/search/identifier"
	"github.com/iot-my-world/brain/security/permission/api"
	"github.com/iot-my-world/brain/security/role"
)

type RecordHandler interface {
	Create(request *CreateRequest) (*CreateResponse, error)
	Retrieve(request *RetrieveRequest) (*RetrieveResponse, error)
	Update(request *UpdateRequest) (*UpdateResponse, error)
}

const ServiceProvider = "Role-RecordHandler"
const CreateService = ServiceProvider + ".Create"
const RetrieveService = ServiceProvider + ".Retrieve"
const UpdateService = ServiceProvider + ".Update"

var SystemUserPermissions = []api.Permission{
	CreateService,
	RetrieveService,
	UpdateService,
}

var CompanyAdminUserPermissions = make([]api.Permission, 0)

var CompanyUserPermissions = make([]api.Permission, 0)

var ClientAdminUserPermissions = make([]api.Permission, 0)

var ClientUserPermissions = make([]api.Permission, 0)

type CreateRequest struct {
	Role role.Role
}

type CreateResponse struct {
	Role role.Role
}

type RetrieveRequest struct {
	Identifier identifier.Identifier
}

type RetrieveResponse struct {
	Role role.Role
}

type UpdateRequest struct {
	Identifier identifier.Identifier
	Role       role.Role
}

type UpdateResponse struct {
	Role role.Role
}

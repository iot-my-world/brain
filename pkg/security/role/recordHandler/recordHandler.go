package recordHandler

import (
	"github.com/iot-my-world/brain/pkg/search/identifier"
	api2 "github.com/iot-my-world/brain/pkg/security/permission/api"
	role2 "github.com/iot-my-world/brain/pkg/security/role"
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

var SystemUserPermissions = []api2.Permission{
	CreateService,
	RetrieveService,
	UpdateService,
}

var CompanyAdminUserPermissions = make([]api2.Permission, 0)

var CompanyUserPermissions = make([]api2.Permission, 0)

var ClientAdminUserPermissions = make([]api2.Permission, 0)

var ClientUserPermissions = make([]api2.Permission, 0)

type CreateRequest struct {
	Role role2.Role
}

type CreateResponse struct {
	Role role2.Role
}

type RetrieveRequest struct {
	Identifier identifier.Identifier
}

type RetrieveResponse struct {
	Role role2.Role
}

type UpdateRequest struct {
	Identifier identifier.Identifier
	Role       role2.Role
}

type UpdateResponse struct {
	Role role2.Role
}

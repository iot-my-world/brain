package administrator

import (
	"github.com/iot-my-world/brain/security/claims"
	"github.com/iot-my-world/brain/security/permission/api"
	apiUser "github.com/iot-my-world/brain/user/api"
)

type Administrator interface {
	Create(request *CreateRequest) (*CreateResponse, error)
	UpdateAllowedFields(request *UpdateAllowedFieldsRequest) (*UpdateAllowedFieldsResponse, error)
}

const ServiceProvider = "APIUser-Administrator"
const CreateService = ServiceProvider + ".Create"
const UpdateAllowedFieldsService = ServiceProvider + ".UpdateAllowedFields"

var SystemUserPermissions = []api.Permission{
	CreateService,
}

var CompanyAdminUserPermissions = make([]api.Permission, 0)

var CompanyUserPermissions = make([]api.Permission, 0)

var ClientAdminUserPermissions = make([]api.Permission, 0)

type CreateRequest struct {
	Claims claims.Claims
	User   apiUser.User
}

type CreateResponse struct {
	User     apiUser.User
	Password string
}

type UpdateAllowedFieldsRequest struct {
	Claims claims.Claims
	User   apiUser.User
}

type UpdateAllowedFieldsResponse struct {
	User apiUser.User
}

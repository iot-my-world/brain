package administrator

import (
	"github.com/iot-my-world/brain/party/client"
	"github.com/iot-my-world/brain/security/claims"
	"github.com/iot-my-world/brain/security/permission/api"
)

type Administrator interface {
	UpdateAllowedFields(request *UpdateAllowedFieldsRequest) (*UpdateAllowedFieldsResponse, error)
	Create(request *CreateRequest) (*CreateResponse, error)
}

const ServiceProvider = "Client-Administrator"
const UpdateAllowedFieldsService = ServiceProvider + ".UpdateAllowedFields"
const CreateService = ServiceProvider + ".Create"

var SystemUserPermissions = make([]api.Permission, 0)

var CompanyAdminUserPermissions = []api.Permission{
	UpdateAllowedFieldsService,
	CreateService,
}

var CompanyUserPermissions = make([]api.Permission, 0)

var ClientAdminUserPermissions = []api.Permission{
	UpdateAllowedFieldsService,
}

var ClientUserPermissions = make([]api.Permission, 0)

type CreateRequest struct {
	Claims claims.Claims
	Client client.Client
}

type CreateResponse struct {
	Client client.Client
}

type UpdateAllowedFieldsRequest struct {
	Claims claims.Claims
	Client client.Client
}

type UpdateAllowedFieldsResponse struct {
	Client client.Client
}

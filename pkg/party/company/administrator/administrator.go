package administrator

import (
	company2 "github.com/iot-my-world/brain/pkg/party/company"
	"github.com/iot-my-world/brain/search/identifier"
	"github.com/iot-my-world/brain/security/claims"
	"github.com/iot-my-world/brain/security/permission/api"
)

type Administrator interface {
	UpdateAllowedFields(request *UpdateAllowedFieldsRequest) (*UpdateAllowedFieldsResponse, error)
	Create(request *CreateRequest) (*CreateResponse, error)
	Delete(request *DeleteRequest) (*DeleteResponse, error)
}

const ServiceProvider = "Company-Administrator"
const UpdateAllowedFieldsService = ServiceProvider + ".UpdateAllowedFields"
const CreateService = ServiceProvider + ".Create"
const DeleteService = ServiceProvider + ".Delete"

var SystemUserPermissions = []api.Permission{
	CreateService,
	DeleteService,
}

var CompanyAdminUserPermissions = []api.Permission{
	UpdateAllowedFieldsService,
}

var CompanyUserPermissions = make([]api.Permission, 0)

var ClientAdminUserPermissions = make([]api.Permission, 0)

var ClientUserPermissions = make([]api.Permission, 0)

type UpdateAllowedFieldsRequest struct {
	Claims  claims.Claims
	Company company2.Company
}

type UpdateAllowedFieldsResponse struct {
	Company company2.Company
}

type CreateRequest struct {
	Claims  claims.Claims
	Company company2.Company
}

type CreateResponse struct {
	Company company2.Company
}

type DeleteRequest struct {
	Claims            claims.Claims
	CompanyIdentifier identifier.Identifier
}

type DeleteResponse struct {
}

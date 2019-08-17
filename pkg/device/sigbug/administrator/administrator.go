package administrator

import (
	"github.com/iot-my-world/brain/pkg/device/sigbug"
	"github.com/iot-my-world/brain/pkg/search/identifier"
	"github.com/iot-my-world/brain/pkg/security/claims"
	"github.com/iot-my-world/brain/pkg/security/permission/api"
	sigfoxBackendDataCallbackMessage "github.com/iot-my-world/brain/pkg/sigfox/backend/callback/data/message"
)

type Administrator interface {
	Create(request *CreateRequest) (*CreateResponse, error)
	UpdateAllowedFields(request *UpdateAllowedFieldsRequest) (*UpdateAllowedFieldsResponse, error)
	LastMessageUpdate(request *LastMessageUpdateRequest) (*LastMessageUpdateResponse, error)
}

const ServiceProvider = "SigbugGPSReading-Administrator"
const UpdateAllowedFieldsService = ServiceProvider + ".UpdateAllowedFields"
const CreateService = ServiceProvider + ".Create"

var SystemUserPermissions = []api.Permission{
	CreateService,
	UpdateAllowedFieldsService,
}

var CompanyAdminUserPermissions = []api.Permission{
	CreateService,
	UpdateAllowedFieldsService,
}

var CompanyUserPermissions = make([]api.Permission, 0)

var ClientAdminUserPermissions = []api.Permission{
	CreateService,
	UpdateAllowedFieldsService,
}

var ClientUserPermissions = make([]api.Permission, 0)

type CreateRequest struct {
	Claims claims.Claims
	Sigbug sigbug.Sigbug
}

type CreateResponse struct {
	Sigbug sigbug.Sigbug
}

type UpdateAllowedFieldsRequest struct {
	Claims claims.Claims
	Sigbug sigbug.Sigbug
}

type UpdateAllowedFieldsResponse struct {
	Sigbug sigbug.Sigbug
}

type LastMessageUpdateRequest struct {
	Claims     claims.Claims
	Identifier identifier.Identifier
	Message    sigfoxBackendDataCallbackMessage.Message
}

type LastMessageUpdateResponse struct {
	Sigbug sigbug.Sigbug
}

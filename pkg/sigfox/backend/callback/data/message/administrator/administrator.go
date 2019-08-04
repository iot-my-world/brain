package administrator

import (
	"github.com/iot-my-world/brain/pkg/search/identifier"
	"github.com/iot-my-world/brain/pkg/security/claims"
	"github.com/iot-my-world/brain/pkg/security/permission/api"
	sigfoxBackendDataCallbackMessage "github.com/iot-my-world/brain/pkg/sigfox/backend/callback/data/message"
)

type Administrator interface {
	Create(request *CreateRequest) (*CreateResponse, error)
}

const ServiceProvider = "MessageDevice-Administrator"
const CreateService = ServiceProvider + ".Create"

var SystemUserPermissions = []api.Permission{
	CreateService,
}

var CompanyAdminUserPermissions = []api.Permission{
	CreateService,
}

var CompanyUserPermissions = make([]api.Permission, 0)

var ClientAdminUserPermissions = []api.Permission{
	CreateService,
}

var ClientUserPermissions = make([]api.Permission, 0)

type CreateRequest struct {
	Claims  claims.Claims
	Message sigfoxBackendDataCallbackMessage.Message
}

type CreateResponse struct {
	Message sigfoxBackendDataCallbackMessage.Message
}

type UpdateAllowedFieldsRequest struct {
	Claims  claims.Claims
	Message sigfoxBackendDataCallbackMessage.Message
}

type UpdateAllowedFieldsResponse struct {
	Message sigfoxBackendDataCallbackMessage.Message
}

type HeartbeatRequest struct {
	Claims            claims.Claims
	MessageIdentifier identifier.Identifier
}

type HeartbeatResponse struct {
}

package recordHandler

import (
	"github.com/iot-my-world/brain/pkg/search/criterion"
	"github.com/iot-my-world/brain/pkg/search/identifier"
	"github.com/iot-my-world/brain/pkg/search/query"
	"github.com/iot-my-world/brain/pkg/security/claims"
	"github.com/iot-my-world/brain/pkg/security/permission/api"
	sigfoxBackendDataCallbackMessage "github.com/iot-my-world/brain/pkg/sigfox/backend/callback/data/message"
)

type RecordHandler interface {
	Create(*CreateRequest) (*CreateResponse, error)
	Retrieve(*RetrieveRequest) (*RetrieveResponse, error)
	Update(*UpdateRequest) (*UpdateResponse, error)
	Delete(*DeleteRequest) (*DeleteResponse, error)
	Collect(*CollectRequest) (*CollectResponse, error)
}

const ServiceProvider = "Message-RecordHandler"
const CreateService = ServiceProvider + ".Create"
const RetrieveService = ServiceProvider + ".Retrieve"
const UpdateService = ServiceProvider + ".Update"
const DeleteService = ServiceProvider + ".Delete"
const CollectService = ServiceProvider + ".Collect"

var SystemUserPermissions = make([]api.Permission, 0)

var CompanyAdminUserPermissions = []api.Permission{
	CollectService,
	RetrieveService,
}

var CompanyUserPermissions = []api.Permission{
	CollectService,
	RetrieveService,
}

var ClientAdminUserPermissions = []api.Permission{
	CollectService,
	RetrieveService,
}

var ClientUserPermissions = []api.Permission{
	CollectService,
	RetrieveService,
}

type CreateRequest struct {
	Message sigfoxBackendDataCallbackMessage.Message
}

type CreateResponse struct {
	Message sigfoxBackendDataCallbackMessage.Message
}

type RetrieveRequest struct {
	Claims     claims.Claims
	Identifier identifier.Identifier
}

type RetrieveResponse struct {
	Message sigfoxBackendDataCallbackMessage.Message
}

type UpdateRequest struct {
	Claims     claims.Claims
	Identifier identifier.Identifier
	Message    sigfoxBackendDataCallbackMessage.Message
}

type UpdateResponse struct{}

type DeleteRequest struct {
	Claims     claims.Claims
	Identifier identifier.Identifier
}

type DeleteResponse struct {
}

type CollectRequest struct {
	Claims   claims.Claims
	Criteria []criterion.Criterion
	Query    query.Query
}

type CollectResponse struct {
	Records []sigfoxBackendDataCallbackMessage.Message
	Total   int
}

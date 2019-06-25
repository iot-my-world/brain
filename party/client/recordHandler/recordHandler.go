package recordHandler

import (
	"github.com/iot-my-world/brain/party/client"
	"github.com/iot-my-world/brain/search/criterion"
	"github.com/iot-my-world/brain/search/identifier"
	"github.com/iot-my-world/brain/search/query"
	"github.com/iot-my-world/brain/security/claims"
	"github.com/iot-my-world/brain/security/permission/api"
)

type RecordHandler interface {
	Create(*CreateRequest) (*CreateResponse, error)
	Retrieve(*RetrieveRequest) (*RetrieveResponse, error)
	Update(*UpdateRequest) (*UpdateResponse, error)
	Delete(*DeleteRequest) (*DeleteResponse, error)
	Collect(*CollectRequest) (*CollectResponse, error)
}

const ServiceProvider = "Client-RecordHandler"
const CreateService = ServiceProvider + ".Create"
const RetrieveService = ServiceProvider + ".Retrieve"
const UpdateService = ServiceProvider + ".Update"
const DeleteService = ServiceProvider + ".Delete"
const CollectService = ServiceProvider + ".Collect"

var CompanyAdminUserPermissions = []api.Permission{
	CreateService,
	RetrieveService,
	UpdateService,
	DeleteService,
	CollectService,
}

var CompanyUserPermissions = []api.Permission{
	CreateService,
	RetrieveService,
	UpdateService,
	DeleteService,
	CollectService,
}

var ClientAdminUserPermissions = []api.Permission{
	CreateService,
	RetrieveService,
	UpdateService,
	DeleteService,
	CollectService,
}

var ClientUserPermissions = []api.Permission{
	CreateService,
	RetrieveService,
	UpdateService,
	DeleteService,
	CollectService,
}

type CreateRequest struct {
	Client client.Client
}

type CreateResponse struct {
	Client client.Client
}

type RetrieveRequest struct {
	Claims     claims.Claims
	Identifier identifier.Identifier
}

type RetrieveResponse struct {
	Client client.Client
}

type UpdateRequest struct {
	Claims     claims.Claims
	Identifier identifier.Identifier
	Client     client.Client
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
	Records []client.Client
	Total   int
}

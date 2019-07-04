package recordHandler

import (
	"github.com/iot-my-world/brain/pkg/search/criterion"
	"github.com/iot-my-world/brain/pkg/search/identifier"
	"github.com/iot-my-world/brain/pkg/search/query"
	sf0012 "github.com/iot-my-world/brain/pkg/tracker/sf001"
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

const ServiceProvider = "SF001-RecordHandler"
const CreateService = ServiceProvider + ".Create"
const RetrieveService = ServiceProvider + ".Retrieve"
const UpdateService = ServiceProvider + ".Update"
const DeleteService = ServiceProvider + ".Delete"
const CollectService = ServiceProvider + ".Collect"

var SystemUserPermissions = make([]api.Permission, 0)

var CompanyAdminUserPermissions = []api.Permission{
	CollectService,
}

var CompanyUserPermissions = []api.Permission{
	CollectService,
}

var ClientAdminUserPermissions = []api.Permission{
	CollectService,
}

var ClientUserPermissions = []api.Permission{
	CollectService,
}

type CreateRequest struct {
	SF001 sf0012.SF001
}

type CreateResponse struct {
	SF001 sf0012.SF001
}

type RetrieveRequest struct {
	Claims     claims.Claims
	Identifier identifier.Identifier
}

type RetrieveResponse struct {
	SF001 sf0012.SF001
}

type UpdateRequest struct {
	Claims     claims.Claims
	Identifier identifier.Identifier
	SF001      sf0012.SF001
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
	Records []sf0012.SF001
	Total   int
}

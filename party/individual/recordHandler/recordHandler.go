package recordHandler

import (
	"github.com/iot-my-world/brain/party/individual"
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

const ServiceProvider = "Individual-RecordHandler"
const CreateService = ServiceProvider + ".Create"
const RetrieveService = ServiceProvider + ".Retrieve"
const UpdateService = ServiceProvider + ".Update"
const DeleteService = ServiceProvider + ".Delete"
const CollectService = ServiceProvider + ".Collect"

var SystemUserPermissions = make([]api.Permission, 0)

var CompanyAdminUserPermissions = make([]api.Permission, 0)

var CompanyUserPermissions = make([]api.Permission, 0)

var ClientAdminUserPermissions = make([]api.Permission, 0)

var ClientUserPermissions = make([]api.Permission, 0)

type CreateRequest struct {
	Individual individual.Individual
}

type CreateResponse struct {
	Individual individual.Individual
}

type RetrieveRequest struct {
	Claims     claims.Claims
	Identifier identifier.Identifier
}

type RetrieveResponse struct {
	Individual individual.Individual
}

type UpdateRequest struct {
	Claims     claims.Claims
	Identifier identifier.Identifier
	Individual individual.Individual
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
	Records []individual.Individual
	Total   int
}

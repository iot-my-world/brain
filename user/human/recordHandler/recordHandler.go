package recordHandler

import (
	"github.com/iot-my-world/brain/search/criterion"
	"github.com/iot-my-world/brain/search/identifier"
	"github.com/iot-my-world/brain/search/query"
	"github.com/iot-my-world/brain/security/claims"
	"github.com/iot-my-world/brain/security/permission/api"
	humanUser "github.com/iot-my-world/brain/user/human"
)

type RecordHandler interface {
	Create(request *CreateRequest) (*CreateResponse, error)
	Retrieve(request *RetrieveRequest) (*RetrieveResponse, error)
	Update(request *UpdateRequest) (*UpdateResponse, error)
	Delete(request *DeleteRequest) (*DeleteResponse, error)
	Collect(request *CollectRequest) (*CollectResponse, error)
}

const ServiceProvider = "HumanUser-RecordHandler"
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

var CompanyUserPermissions = make([]api.Permission, 0)

var ClientAdminUserPermissions = []api.Permission{
	CollectService,
	RetrieveService,
}

var ClientUserPermissions = make([]api.Permission, 0)

type CreateRequest struct {
	User humanUser.User
}

type CreateResponse struct {
	User humanUser.User
}

type DeleteRequest struct {
	Claims     claims.Claims
	Identifier identifier.Identifier
}

type DeleteResponse struct {
	User humanUser.User
}

type UpdateRequest struct {
	Claims     claims.Claims
	Identifier identifier.Identifier
	User       humanUser.User
}

type UpdateResponse struct {
	User humanUser.User
}

type RetrieveRequest struct {
	Claims     claims.Claims
	Identifier identifier.Identifier
}

type RetrieveResponse struct {
	User humanUser.User
}

type CollectRequest struct {
	Claims   claims.Claims
	Criteria []criterion.Criterion
	Query    query.Query
}

type CollectResponse struct {
	Records []humanUser.User
	Total   int
}

package recordHandler

import (
	"github.com/iot-my-world/brain/pkg/search/criterion"
	"github.com/iot-my-world/brain/pkg/search/identifier"
	"github.com/iot-my-world/brain/pkg/search/query"
	"github.com/iot-my-world/brain/pkg/user/human"
	"github.com/iot-my-world/brain/security/claims"
	"github.com/iot-my-world/brain/security/permission/api"
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
	User human.User
}

type CreateResponse struct {
	User human.User
}

type DeleteRequest struct {
	Claims     claims.Claims
	Identifier identifier.Identifier
}

type DeleteResponse struct {
	User human.User
}

type UpdateRequest struct {
	Claims     claims.Claims
	Identifier identifier.Identifier
	User       human.User
}

type UpdateResponse struct {
	User human.User
}

type RetrieveRequest struct {
	Claims     claims.Claims
	Identifier identifier.Identifier
}

type RetrieveResponse struct {
	User human.User
}

type CollectRequest struct {
	Claims   claims.Claims
	Criteria []criterion.Criterion
	Query    query.Query
}

type CollectResponse struct {
	Records []human.User
	Total   int
}

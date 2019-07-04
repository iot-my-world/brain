package recordHandler

import (
	"github.com/iot-my-world/brain/api"
	system2 "github.com/iot-my-world/brain/pkg/party/system"
	"github.com/iot-my-world/brain/search/criterion"
	"github.com/iot-my-world/brain/search/identifier"
	"github.com/iot-my-world/brain/search/query"
	"github.com/iot-my-world/brain/security/claims"
	apiPermission "github.com/iot-my-world/brain/security/permission/api"
	"github.com/iot-my-world/brain/validate/reasonInvalid"
)

type RecordHandler interface {
	Create(request *CreateRequest) (*CreateResponse, error)
	Retrieve(request *RetrieveRequest) (*RetrieveResponse, error)
	Update(request *UpdateRequest) (*UpdateResponse, error)
	Delete(request *DeleteRequest) (*DeleteResponse, error)
	Validate(request *ValidateRequest) (*ValidateResponse, error)
	Collect(request *CollectRequest) (*CollectResponse, error)
}

const Create api.Method = "Create"
const Retrieve api.Method = "Retrieve"
const Update api.Method = "Update"
const Delete api.Method = "Delete"
const Validate api.Method = "Validate"

const ServiceProvider = "System-RecordHandler"
const CreateService = ServiceProvider + ".Create"
const RetrieveService = ServiceProvider + ".Retrieve"
const UpdateService = ServiceProvider + ".Update"
const DeleteService = ServiceProvider + ".Delete"
const CollectService = ServiceProvider + ".Collect"

var SystemUserPermissions = make([]apiPermission.Permission, 0)

var CompanyAdminUserPermissions = []apiPermission.Permission{
	CollectService,
}

var CompanyUserPermissions = make([]apiPermission.Permission, 0)

var ClientAdminUserPermissions = []apiPermission.Permission{
	CollectService,
}

var ClientUserPermissions = make([]apiPermission.Permission, 0)

type CollectRequest struct {
	Claims   claims.Claims
	Criteria []criterion.Criterion
	Query    query.Query
}

type CollectResponse struct {
	Records []system2.System
	Total   int
}

type ValidateRequest struct {
	System system2.System
	Method api.Method
}

type ValidateResponse struct {
	ReasonsInvalid []reasonInvalid.ReasonInvalid
}

type CreateRequest struct {
	System system2.System
}

type CreateResponse struct {
	System system2.System
}

type DeleteRequest struct {
	Claims     claims.Claims
	Identifier identifier.Identifier
}

type DeleteResponse struct {
	System system2.System
}

type UpdateRequest struct {
	Claims     claims.Claims
	Identifier identifier.Identifier
	System     system2.System
}

type UpdateResponse struct {
	System system2.System
}

type RetrieveRequest struct {
	Claims     claims.Claims
	Identifier identifier.Identifier
}

type RetrieveResponse struct {
	System system2.System
}

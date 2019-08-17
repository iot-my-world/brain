package recordHandler

import (
	sigbugGPSReading "github.com/iot-my-world/brain/pkg/device/sigbug/reading/gps"
	"github.com/iot-my-world/brain/pkg/search/criterion"
	"github.com/iot-my-world/brain/pkg/search/identifier"
	"github.com/iot-my-world/brain/pkg/search/query"
	"github.com/iot-my-world/brain/pkg/security/claims"
	"github.com/iot-my-world/brain/pkg/security/permission/api"
)

type RecordHandler interface {
	Create(*CreateRequest) (*CreateResponse, error)
	Retrieve(*RetrieveRequest) (*RetrieveResponse, error)
	Update(*UpdateRequest) (*UpdateResponse, error)
	Delete(*DeleteRequest) (*DeleteResponse, error)
	Collect(*CollectRequest) (*CollectResponse, error)
}

const ServiceProvider = "Reading-RecordHandler"
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
	Reading sigbugGPSReading.Reading
}

type CreateResponse struct {
	Reading sigbugGPSReading.Reading
}

type RetrieveRequest struct {
	Claims     claims.Claims
	Identifier identifier.Identifier
}

type RetrieveResponse struct {
	Reading sigbugGPSReading.Reading
}

type UpdateRequest struct {
	Claims     claims.Claims
	Identifier identifier.Identifier
	Reading    sigbugGPSReading.Reading
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
	Records []sigbugGPSReading.Reading
	Total   int
}

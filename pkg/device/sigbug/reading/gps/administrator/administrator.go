package administrator

import (
	sigbugGPSReading "github.com/iot-my-world/brain/pkg/device/sigbug/reading/gps"
	"github.com/iot-my-world/brain/pkg/security/claims"
	"github.com/iot-my-world/brain/pkg/security/permission/api"
)

type Administrator interface {
	Create(request *CreateRequest) (*CreateResponse, error)
}

const ServiceProvider = "SigbugGPSReading-Administrator"
const CreateService = ServiceProvider + ".Create"

var SystemUserPermissions = []api.Permission{
	CreateService,
}

var CompanyAdminUserPermissions = make([]api.Permission, 0)

var CompanyUserPermissions = make([]api.Permission, 0)

var ClientAdminUserPermissions = make([]api.Permission, 0)

var ClientUserPermissions = make([]api.Permission, 0)

type CreateRequest struct {
	Claims  claims.Claims
	Reading sigbugGPSReading.Reading
}

type CreateResponse struct {
	Reading sigbugGPSReading.Reading
}

type UpdateAllowedFieldsRequest struct {
	Claims  claims.Claims
	Reading sigbugGPSReading.Reading
}

type UpdateAllowedFieldsResponse struct {
	Reading sigbugGPSReading.Reading
}

package administrator

import (
	"github.com/iot-my-world/brain/security/claims"
	apiUser "github.com/iot-my-world/brain/user/api"
)

type Administrator interface {
	Create(request *CreateRequest) (*CreateResponse, error)
	UpdateAllowedFields(request *UpdateAllowedFieldsRequest) (*UpdateAllowedFieldsResponse, error)
}

type CreateRequest struct {
	Claims claims.Claims
	User   apiUser.User
}

type CreateResponse struct {
	User     apiUser.User
	Password string
}

type UpdateAllowedFieldsRequest struct {
	Claims claims.Claims
	User   apiUser.User
}

type UpdateAllowedFieldsResponse struct {
	User apiUser.User
}

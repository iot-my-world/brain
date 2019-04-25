package administrator

import (
	"gitlab.com/iotTracker/brain/security/claims"
	apiUser "gitlab.com/iotTracker/brain/user/api"
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

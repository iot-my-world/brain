package administrator

import (
	"gitlab.com/iotTracker/brain/party/user"
	"gitlab.com/iotTracker/brain/security/claims"
)

type Administrator interface {
	GetMyUser(request *GetMyUserRequest, response *GetMyUserResponse) error
	UpdateAllowedFields(request *UpdateAllowedFieldsRequest, response *UpdateAllowedFieldsResponse) error
}

type UpdateAllowedFieldsRequest struct {
	Claims claims.Claims
	User   user.User
}

type UpdateAllowedFieldsResponse struct {
	User user.User
}

// GetMyUserRequest is the request object for the Handler GetMyUser service
type GetMyUserRequest struct {
	Claims claims.Claims
}

// GetMyUserResponse is the response object for the Handler GetMyUser service
type GetMyUserResponse struct {
	User user.User
}

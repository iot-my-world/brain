package administrator

import (
	"gitlab.com/iotTracker/brain/search/identifier"
	"gitlab.com/iotTracker/brain/security/claims"
	"gitlab.com/iotTracker/brain/user"
)

type Administrator interface {
	GetMyUser(request *GetMyUserRequest) (*GetMyUserResponse, error)
	UpdateAllowedFields(request *UpdateAllowedFieldsRequest) (*UpdateAllowedFieldsResponse, error)
	Create(request *CreateRequest) (*CreateResponse, error)
	ChangePassword(request *ChangePasswordRequest) (*ChangePasswordResponse, error)
}

type UpdateAllowedFieldsRequest struct {
	Claims claims.Claims
	User   user.User
}

type UpdateAllowedFieldsResponse struct {
	User user.User
}

type GetMyUserRequest struct {
	Claims claims.Claims
}

type GetMyUserResponse struct {
	User user.User
}

type CreateRequest struct {
	Claims claims.Claims
	User   user.User
}

type CreateResponse struct {
	User user.User
}

type ChangePasswordRequest struct {
	Claims      claims.Claims
	Identifier  identifier.Identifier
	NewPassword string
}

type ChangePasswordResponse struct {
	User user.User
}

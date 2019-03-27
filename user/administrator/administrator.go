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
	SetPassword(request *SetPasswordRequest) (*SetPasswordResponse, error)
	CheckPassword(request *CheckPasswordRequest) (*CheckPasswordResponse, error)
	UpdatePassword(request *UpdatePasswordRequest) (*UpdatePasswordResponse, error)
	ForgotPassword(request *ForgotPasswordRequest) (*ForgotPasswordResponse, error)
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

type SetPasswordRequest struct {
	Claims      claims.Claims
	Identifier  identifier.Identifier
	NewPassword string
}

type SetPasswordResponse struct {
	User user.User
}

type CheckPasswordRequest struct {
	Claims   claims.Claims
	Password string
}

type CheckPasswordResponse struct {
	Result bool
}

type UpdatePasswordRequest struct {
	Claims           claims.Claims
	ExistingPassword string
	NewPassword      string
}

type UpdatePasswordResponse struct {
	User user.User
}

type ForgotPasswordRequest struct {
	UsernameOrEmailAddress string
}

type ForgotPasswordResponse struct {
}

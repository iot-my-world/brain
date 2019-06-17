package administrator

import (
	"gitlab.com/iotTracker/brain/search/identifier"
	"gitlab.com/iotTracker/brain/security/claims"
	humanUser "gitlab.com/iotTracker/brain/user/human"
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
	User   humanUser.User
}

type UpdateAllowedFieldsResponse struct {
	User humanUser.User
}

type GetMyUserRequest struct {
	Claims claims.Claims
}

type GetMyUserResponse struct {
	User humanUser.User
}

type CreateRequest struct {
	Claims claims.Claims
	User   humanUser.User
}

type CreateResponse struct {
	User humanUser.User
}

type SetPasswordRequest struct {
	Claims      claims.Claims
	Identifier  identifier.Identifier
	NewPassword string
}

type SetPasswordResponse struct {
	User humanUser.User
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
	User humanUser.User
}

type ForgotPasswordRequest struct {
	UsernameOrEmailAddress string
}

type ForgotPasswordResponse struct {
}

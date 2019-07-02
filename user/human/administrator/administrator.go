package administrator

import (
	"github.com/iot-my-world/brain/search/identifier"
	"github.com/iot-my-world/brain/security/claims"
	"github.com/iot-my-world/brain/security/permission/api"
	humanUser "github.com/iot-my-world/brain/user/human"
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

const ServiceProvider = "HumanUser-Administrator"
const GetMyUserService = ServiceProvider + ".GetMyUser"
const UpdateAllowedFieldsService = ServiceProvider + ".UpdateAllowedFields"
const CreateService = ServiceProvider + ".Create"
const SetPasswordService = ServiceProvider + ".SetPassword"
const CheckPasswordService = ServiceProvider + ".CheckPassword"
const UpdatePasswordService = ServiceProvider + ".UpdatePassword"
const ForgotPasswordService = ServiceProvider + ".ForgotPassword"

var SystemUserPermissions = make([]api.Permission, 0)

var CompanyAdminUserPermissions = []api.Permission{
	UpdateAllowedFieldsService,
	CreateService,
	GetMyUserService,
	UpdatePasswordService,
	CheckPasswordService,
}

var CompanyUserPermissions = []api.Permission{
	UpdateAllowedFieldsService,
	CreateService,
	GetMyUserService,
	UpdatePasswordService,
	CheckPasswordService,
}

var ClientAdminUserPermissions = []api.Permission{
	UpdateAllowedFieldsService,
	CreateService,
	GetMyUserService,
	UpdatePasswordService,
	CheckPasswordService,
}

var ClientUserPermissions = []api.Permission{
	UpdateAllowedFieldsService,
	CreateService,
	GetMyUserService,
	UpdatePasswordService,
	CheckPasswordService,
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
	URLToken string
}

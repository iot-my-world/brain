package administrator

import (
	"github.com/iot-my-world/brain/pkg/search/identifier"
	"github.com/iot-my-world/brain/pkg/security/claims"
	"github.com/iot-my-world/brain/pkg/security/permission/api"
	"github.com/iot-my-world/brain/pkg/user/human"
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
	User   human.User
}

type UpdateAllowedFieldsResponse struct {
	User human.User
}

type GetMyUserRequest struct {
	Claims claims.Claims
}

type GetMyUserResponse struct {
	User human.User
}

type CreateRequest struct {
	Claims claims.Claims
	User   human.User
}

type CreateResponse struct {
	User human.User
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
	User human.User
}

type ForgotPasswordRequest struct {
	UsernameOrEmailAddress string
}

type ForgotPasswordResponse struct {
	URLToken string
}

package registrar

import (
	"gitlab.com/iotTracker/brain/party/user"
	"gitlab.com/iotTracker/brain/security/claims"
	"gitlab.com/iotTracker/brain/api"
)

type Registrar interface {
	RegisterSystemAdminUser(request *RegisterSystemAdminUserRequest, response *RegisterSystemAdminUserResponse) error

	InviteCompanyAdminUser(request *InviteCompanyAdminUserRequest, response *InviteCompanyAdminUserResponse) error
	RegisterCompanyAdminUser(request *RegisterCompanyAdminUserRequest, response *RegisterCompanyAdminUserResponse) error
	InviteCompanyUser(request *InviteCompanyUserRequest, response *InviteCompanyUserResponse) error

	InviteClientAdminUser(request *InviteClientAdminUserRequest, response *InviteClientAdminUserResponse) error
	RegisterClientAdminUser(request *RegisterClientAdminUserRequest, response *RegisterClientAdminUserResponse) error
}

const InviteCompanyAdminUser api.Method = "InviteCompanyAdminUser"
const RegisterCompanyAdminUser api.Method = "RegisterCompanyAdminUser"
const InviteCompanyUser api.Method = "InviteCompanyUser"
const RegisterCompanyUser api.Method = "RegisterCompanyUser"

const InviteClientAdminUser api.Method = "InviteClientAdminUser"
const RegisterClientAdminUser api.Method = "RegisterClientAdminUser"
const InviteClientUser api.Method = "InviteClientUser"
const RegisterClientUser api.Method = "RegisterClientUser"

type RegisterSystemAdminUserRequest struct {
	Claims   claims.Claims
	User     user.User
	Password string
}

type RegisterSystemAdminUserResponse struct {
	User user.User
}

type InviteCompanyAdminUserRequest struct {
	// claims for company party retrieval
	Claims claims.Claims
	// the minimal company admin user
	User user.User
}

type InviteCompanyAdminUserResponse struct {
	URLToken string
}

type RegisterCompanyAdminUserRequest struct {
	Claims   claims.Claims
	User     user.User
	Password string
}

type RegisterCompanyAdminUserResponse struct {
	User user.User
}

type InviteCompanyUserRequest struct {
	// claims for company party retrieval
	Claims claims.Claims
	// the minimal company user
	User user.User
}

type InviteCompanyUserResponse struct {
	URLToken string
}

type InviteClientAdminUserRequest struct {
	// claims for client party retrieval
	Claims claims.Claims
	// the minimal client admin user
	User user.User
}

type InviteClientAdminUserResponse struct {
	URLToken string
}

type RegisterClientAdminUserRequest struct {
	Claims   claims.Claims
	User     user.User
	Password string
}

type RegisterClientAdminUserResponse struct {
	User user.User
}

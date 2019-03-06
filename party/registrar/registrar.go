package registrar

import (
	"gitlab.com/iotTracker/brain/party/user"
	"gitlab.com/iotTracker/brain/search/identifier"
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

const Invite api.Method = "Invite"

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
	// an identifier to retrieve the company
	// to which the admin user will belong
	CompanyIdentifier identifier.Identifier
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
	// an identifier to retrieve the company
	// to which the new user will belong
	CompanyIdentifier identifier.Identifier
	// Address to which the registration
	// invite will be sent
	EmailAddress string
}

type InviteCompanyUserResponse struct {
	URLToken string
}

type InviteClientAdminUserRequest struct {
	// claims for client party retrieval
	Claims          claims.Claims
	// an identifier to retrieve the client
	// to which the admin user will belong
	ClientIdentifier identifier.Identifier
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

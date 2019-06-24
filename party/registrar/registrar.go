package registrar

import (
	"github.com/iot-my-world/brain/search/identifier"
	"github.com/iot-my-world/brain/search/identifier/party"
	"github.com/iot-my-world/brain/security/claims"
	humanUser "github.com/iot-my-world/brain/user/human"
)

type Registrar interface {
	RegisterSystemAdminUser(request *RegisterSystemAdminUserRequest) (*RegisterSystemAdminUserResponse, error)

	InviteCompanyAdminUser(request *InviteCompanyAdminUserRequest) (*InviteCompanyAdminUserResponse, error)
	RegisterCompanyAdminUser(request *RegisterCompanyAdminUserRequest) (*RegisterCompanyAdminUserResponse, error)
	InviteCompanyUser(request *InviteCompanyUserRequest) (*InviteCompanyUserResponse, error)
	RegisterCompanyUser(request *RegisterCompanyUserRequest) (*RegisterCompanyUserResponse, error)

	InviteClientAdminUser(request *InviteClientAdminUserRequest) (*InviteClientAdminUserResponse, error)
	RegisterClientAdminUser(request *RegisterClientAdminUserRequest) (*RegisterClientAdminUserResponse, error)
	InviteClientUser(request *InviteClientUserRequest) (*InviteClientUserResponse, error)
	RegisterClientUser(request *RegisterClientUserRequest) (*RegisterClientUserResponse, error)

	InviteUser(request *InviteUserRequest) (*InviteUserResponse, error)

	AreAdminsRegistered(request *AreAdminsRegisteredRequest) (*AreAdminsRegisteredResponse, error)
}

type RegisterSystemAdminUserRequest struct {
	Claims claims.Claims
	User   humanUser.User
}

type RegisterSystemAdminUserResponse struct {
	User humanUser.User
}

type InviteCompanyAdminUserRequest struct {
	Claims            claims.Claims
	CompanyIdentifier identifier.Identifier
}

type InviteCompanyAdminUserResponse struct {
	URLToken string
}

type RegisterCompanyAdminUserRequest struct {
	Claims claims.Claims
	User   humanUser.User
}

type RegisterCompanyAdminUserResponse struct{}

type InviteCompanyUserRequest struct {
	Claims         claims.Claims
	UserIdentifier identifier.Identifier
}

type InviteCompanyUserResponse struct {
	URLToken string
}

type RegisterCompanyUserRequest struct {
	Claims claims.Claims
	User   humanUser.User
}

type RegisterCompanyUserResponse struct{}

type InviteClientAdminUserRequest struct {
	Claims           claims.Claims
	ClientIdentifier identifier.Identifier
}

type InviteClientAdminUserResponse struct {
	URLToken string
}

type RegisterClientAdminUserRequest struct {
	Claims claims.Claims
	User   humanUser.User
}

type RegisterClientAdminUserResponse struct {
}

type InviteClientUserRequest struct {
	Claims         claims.Claims
	UserIdentifier identifier.Identifier
}

type InviteClientUserResponse struct {
	URLToken string
}

type RegisterClientUserRequest struct {
	Claims claims.Claims
	User   humanUser.User
}

type RegisterClientUserResponse struct {
}

type InviteUserRequest struct {
	Claims         claims.Claims
	UserIdentifier identifier.Identifier
}

type InviteUserResponse struct {
	URLToken string
}

type AreAdminsRegisteredRequest struct {
	Claims           claims.Claims
	PartyIdentifiers []party.Identifier
}

type AreAdminsRegisteredResponse struct {
	Result map[string]bool
}

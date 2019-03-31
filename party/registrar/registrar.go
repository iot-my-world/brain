package registrar

import (
	"gitlab.com/iotTracker/brain/search/identifier"
	"gitlab.com/iotTracker/brain/security/claims"
	"gitlab.com/iotTracker/brain/user"
	"gitlab.com/iotTracker/brain/search/identifier/party"
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
	User   user.User
}

type RegisterSystemAdminUserResponse struct {
	User user.User
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
	User   user.User
}

type RegisterCompanyAdminUserResponse struct {
	User user.User
}

type InviteCompanyUserRequest struct {
	Claims         claims.Claims
	UserIdentifier identifier.Identifier
}

type InviteCompanyUserResponse struct {
	URLToken string
}

type RegisterCompanyUserRequest struct {
	Claims claims.Claims
	User   user.User
}

type RegisterCompanyUserResponse struct {
	User user.User
}

type InviteClientAdminUserRequest struct {
	Claims           claims.Claims
	ClientIdentifier identifier.Identifier
}

type InviteClientAdminUserResponse struct {
	URLToken string
}

type RegisterClientAdminUserRequest struct {
	Claims claims.Claims
	User   user.User
}

type RegisterClientAdminUserResponse struct {
	User user.User
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
	User   user.User
}

type RegisterClientUserResponse struct {
	User user.User
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

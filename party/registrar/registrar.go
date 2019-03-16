package registrar

import (
	"gitlab.com/iotTracker/brain/party"
	"gitlab.com/iotTracker/brain/party/user"
	"gitlab.com/iotTracker/brain/search/identifier"
	"gitlab.com/iotTracker/brain/security/claims"
)

type Registrar interface {
	RegisterSystemAdminUser(request *RegisterSystemAdminUserRequest, response *RegisterSystemAdminUserResponse) error

	InviteCompanyAdminUser(request *InviteCompanyAdminUserRequest, response *InviteCompanyAdminUserResponse) error
	RegisterCompanyAdminUser(request *RegisterCompanyAdminUserRequest, response *RegisterCompanyAdminUserResponse) error
	InviteCompanyUser(request *InviteCompanyUserRequest, response *InviteCompanyUserResponse) error
	RegisterCompanyUser(request *RegisterCompanyUserRequest, response *RegisterCompanyUserResponse) error

	InviteClientAdminUser(request *InviteClientAdminUserRequest, response *InviteClientAdminUserResponse) error
	RegisterClientAdminUser(request *RegisterClientAdminUserRequest, response *RegisterClientAdminUserResponse) error
	InviteClientUser(request *InviteClientUserRequest, response *InviteClientUserResponse) error
	RegisterClientUser(request *RegisterClientUserRequest, response *RegisterClientUserResponse) error

	AreAdminsRegistered(request *AreAdminsRegisteredRequest, response *AreAdminsRegisteredResponse) error
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
	Claims claims.Claims
	User   user.User
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
	Claims claims.Claims
	User   user.User
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

type AreAdminsRegisteredRequest struct {
	Claims       claims.Claims
	PartyDetails []party.Detail
}

type AreAdminsRegisteredResponse struct {
	Result map[string]bool
}

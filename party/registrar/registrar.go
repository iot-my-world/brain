package registrar

import (
	"gitlab.com/iotTracker/brain/party/user"
	"gitlab.com/iotTracker/brain/search/identifier"
	"gitlab.com/iotTracker/brain/security/claims"
)

type Registrar interface {
	InviteCompanyAdminUser(request *InviteCompanyAdminUserRequest, response *InviteCompanyAdminUserResponse) error
	RegisterCompanyAdminUser(request *RegisterCompanyAdminUserRequest, response *RegisterCompanyAdminUserResponse) error
}

type InviteCompanyAdminUserRequest struct {
	PartyIdentifier identifier.Identifier
}

type InviteCompanyAdminUserResponse struct {
}

type RegisterCompanyAdminUserRequest struct {
	Claims   claims.Claims
	User     user.User
	Password string
}

type RegisterCompanyAdminUserResponse struct {
	User user.User
}

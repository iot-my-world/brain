package registrar

import "gitlab.com/iotTracker/brain/search/identifier"

type Registrar interface {
	InviteCompanyAdminUser(request *InviteCompanyAdminUserRequest, response *InviteCompanyAdminUserResponse) error
}

type InviteCompanyAdminUserRequest struct {
	PartyIdentifier identifier.Identifier
}

type InviteCompanyAdminUserResponse struct {
}

type RegisterAdminUserRequest struct {
}

type RegisterAdminUserResponse struct {
}

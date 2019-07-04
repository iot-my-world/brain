package administrator

import (
	"github.com/iot-my-world/brain/pkg/party"
	client2 "github.com/iot-my-world/brain/pkg/party/client"
	company2 "github.com/iot-my-world/brain/pkg/party/company"
	"github.com/iot-my-world/brain/search/identifier"
	"github.com/iot-my-world/brain/security/claims"
	"github.com/iot-my-world/brain/security/permission/api"
)

type Administrator interface {
	GetMyParty(request *GetMyPartyRequest) (*GetMyPartyResponse, error)
	RetrieveParty(request *RetrievePartyRequest) (*RetrievePartyResponse, error)
	CreateAndInviteCompany(request *CreateAndInviteCompanyRequest) (*CreateAndInviteCompanyResponse, error)
	CreateAndInviteClient(request *CreateAndInviteClientRequest) (*CreateAndInviteClientResponse, error)
}

const ServiceProvider = "Party-Administrator"
const GetMyPartyService = ServiceProvider + ".GetMyParty"
const RetrievePartyService = ServiceProvider + ".RetrieveParty"
const CreateAndInviteCompanyService = ServiceProvider + ".CreateAndInviteCompany"
const CreateAndInviteClientService = ServiceProvider + ".CreateAndInviteClient"

var SystemUserPermissions = make([]api.Permission, 0)

var CompanyAdminUserPermissions = []api.Permission{
	GetMyPartyService,
	RetrievePartyService,
}

var CompanyUserPermissions = []api.Permission{
	GetMyPartyService,
	RetrievePartyService,
}

var ClientAdminUserPermissions = []api.Permission{
	GetMyPartyService,
	RetrievePartyService,
}

var ClientUserPermissions = []api.Permission{
	GetMyPartyService,
	RetrievePartyService,
}

type GetMyPartyRequest struct {
	Claims claims.Claims
}

type GetMyPartyResponse struct {
	Party     party.Party
	PartyType party.Type
}

type RetrievePartyRequest struct {
	Claims     claims.Claims
	PartyType  party.Type
	Identifier identifier.Identifier
}

type RetrievePartyResponse struct {
	Party party.Party
}

type CreateAndInviteCompanyRequest struct {
	Company company2.Company
}

type CreateAndInviteCompanyResponse struct {
	RegistrationURLToken string
}

type CreateAndInviteClientRequest struct {
	Client client2.Client
}

type CreateAndInviteClientResponse struct {
	RegistrationURLToken string
}

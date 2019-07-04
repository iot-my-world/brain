package administrator

import (
	"github.com/iot-my-world/brain/party"
	"github.com/iot-my-world/brain/party/client"
	"github.com/iot-my-world/brain/party/company"
	"github.com/iot-my-world/brain/search/identifier"
	"github.com/iot-my-world/brain/security/claims"
	"github.com/iot-my-world/brain/security/permission/api"
)

type Administrator interface {
	GetMyParty(request *GetMyPartyRequest) (*GetMyPartyResponse, error)
	RetrieveParty(request *RetrievePartyRequest) (*RetrievePartyResponse, error)
	CreateAndInviteCompany(request *CreateAndInviteCompanyRequest) (*CreateAndInviteCompanyResponse, error)
	CreateAndInviteCompanyClient(request *CreateAndInviteCompanyClientRequest) (*CreateAndInviteCompanyClientResponse, error)
	CreateAndInviteIndividualClient(request *CreateAndInviteIndividualClientRequest) (*CreateAndInviteIndividualClientResponse, error)
}

const ServiceProvider = "Party-Administrator"
const GetMyPartyService = ServiceProvider + ".GetMyParty"
const RetrievePartyService = ServiceProvider + ".RetrieveParty"
const CreateAndInviteCompanyService = ServiceProvider + ".CreateAndInviteCompany"
const CreateAndInviteCompanyClientService = ServiceProvider + ".CreateAndInviteCompanyClient"
const CreateAndInviteIndividualClientService = ServiceProvider + ".CreateAndInviteIndividualClient"

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
	Company company.Company
}

type CreateAndInviteCompanyResponse struct {
	RegistrationURLToken string
}

type CreateAndInviteCompanyClientRequest struct {
	Client client.Client
}

type CreateAndInviteCompanyClientResponse struct {
	RegistrationURLToken string
}

type CreateAndInviteIndividualClientRequest struct {
	Client client.Client
}

type CreateAndInviteIndividualClientResponse struct {
	RegistrationURLToken string
}

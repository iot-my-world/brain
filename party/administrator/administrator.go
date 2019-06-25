package administrator

import (
	"github.com/iot-my-world/brain/party"
	"github.com/iot-my-world/brain/search/identifier"
	"github.com/iot-my-world/brain/security/claims"
	"github.com/iot-my-world/brain/security/permission/api"
)

type Administrator interface {
	GetMyParty(request *GetMyPartyRequest) (*GetMyPartyResponse, error)
	RetrieveParty(request *RetrievePartyRequest) (*RetrievePartyResponse, error)
}

const ServiceProvider = "Party-Administrator"
const GetMyPartyService = ServiceProvider + ".GetMyParty"
const RetrievePartyService = ServiceProvider + ".RetrieveParty"

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

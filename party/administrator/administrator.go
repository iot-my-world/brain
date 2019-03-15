package administrator

import (
	"gitlab.com/iotTracker/brain/party"
	"gitlab.com/iotTracker/brain/search/identifier"
	"gitlab.com/iotTracker/brain/security/claims"
)

type Administrator interface {
	GetMyParty(request *GetMyPartyRequest, response *GetMyPartyResponse) error
	RetrieveParty(request *RetrievePartyRequest, response *RetrievePartyResponse) error
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
	Party     party.Party
	PartyType party.Type
}

package administrator

import (
	"github.com/iot-my-world/brain/party"
	"github.com/iot-my-world/brain/search/identifier"
	"github.com/iot-my-world/brain/security/claims"
)

type Administrator interface {
	GetMyParty(request *GetMyPartyRequest) (*GetMyPartyResponse, error)
	RetrieveParty(request *RetrievePartyRequest) (*RetrievePartyResponse, error)
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

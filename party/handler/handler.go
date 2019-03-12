package handler

import (
	"gitlab.com/iotTracker/brain/party"
	"gitlab.com/iotTracker/brain/search/identifier"
	"gitlab.com/iotTracker/brain/security/claims"
)

// Handler is a generic way in which to
type Handler interface {
	GetMyParty(request *GetMyPartyRequest, response *GetMyPartyResponse) error
	RetrieveParty(request *RetrievePartyRequest, resposne *RetrievePartyResponse) error
}

type GetMyPartyRequest struct {
	Claims claims.Claims
}

type GetMyPartyResponse struct {
	Party     interface{}
	PartyType party.Type
}

type RetrievePartyRequest struct {
	Claims     claims.Claims
	PartyType  party.Type
	Identifier identifier.Identifier
}

type RetrievePartyResponse struct {
	Party     interface{}
	PartyType party.Type
}

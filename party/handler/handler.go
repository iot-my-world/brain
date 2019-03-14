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

// GetMyPartyRequest is the request object for the Handler GetMyParty service
type GetMyPartyRequest struct {
	Claims claims.Claims
}

// GetMyPartyResponse is the response object for the Handler GetMyParty service
type GetMyPartyResponse struct {
	Party     party.Party
	PartyType party.Type
}

// RetrievePartyRequest is the request object for the Handler RetrieveParty service
type RetrievePartyRequest struct {
	Claims     claims.Claims
	PartyType  party.Type
	Identifier identifier.Identifier
}

// RetrievePartyResponse is the response object for the Handler RetrieveParty service
type RetrievePartyResponse struct {
	Party     party.Party
	PartyType party.Type
}

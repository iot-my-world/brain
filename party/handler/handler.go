package handler

import (
	"gitlab.com/iotTracker/brain/security/claims"
	"gitlab.com/iotTracker/brain/party"
)

type Handler interface {
	GetMyParty(request *GetMyPartyRequest, response *GetMyPartyResponse) error
}

type GetMyPartyRequest struct {
	Claims claims.Claims
}

type GetMyPartyResponse struct {
	Party     interface{}
	PartyType party.Type
}

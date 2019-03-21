package jsonRpc

import (
	"gitlab.com/iotTracker/brain/log"
	"gitlab.com/iotTracker/brain/party"
	partyAdministrator "gitlab.com/iotTracker/brain/party/administrator"
	wrappedParty "gitlab.com/iotTracker/brain/party/wrapped"
	wrappedIdentifier "gitlab.com/iotTracker/brain/search/identifier/wrapped"
	"gitlab.com/iotTracker/brain/security/wrappedClaims"
	"net/http"
)

type adaptor struct {
	partyAdministrator partyAdministrator.Administrator
}

func New(
	partyAdministrator partyAdministrator.Administrator,
) *adaptor {
	return &adaptor{
		partyAdministrator: partyAdministrator,
	}
}

type GetMyPartyRequest struct{}

type GetMyPartyResponse struct {
	Party     interface{} `json:"party"`
	PartyType party.Type  `json:"partyType"`
}

func (a *adaptor) GetMyParty(r *http.Request, request *GetMyPartyRequest, response *GetMyPartyResponse) error {
	claims, err := wrappedClaims.UnwrapClaimsFromContext(r)
	if err != nil {
		log.Warn(err.Error())
		return err
	}

	getMyPartyResponse, err := a.partyAdministrator.GetMyParty(&partyAdministrator.GetMyPartyRequest{
		Claims: claims,
	})
	if err != nil {
		return err
	}

	response.Party = getMyPartyResponse.Party
	response.PartyType = getMyPartyResponse.PartyType

	return nil
}

type RetrievePartyRequest struct {
	PartyType  party.Type
	Identifier wrappedIdentifier.Wrapped
}

type RetrievePartyResponse struct {
	Party wrappedParty.Wrapped `json:"Party"`
}

func (a *adaptor) RetrieveParty(r *http.Request, request *RetrievePartyRequest, response *RetrievePartyResponse) error {
	claims, err := wrappedClaims.UnwrapClaimsFromContext(r)
	if err != nil {
		log.Warn(err.Error())
		return err
	}

	partyIdentifier, err := request.Identifier.UnWrap()
	if err != nil {
		return err
	}

	retrievePartyResponse, err := a.partyAdministrator.RetrieveParty(&partyAdministrator.RetrievePartyRequest{
		Claims:     claims,
		PartyType:  request.PartyType,
		Identifier: partyIdentifier,
	})
	if err != nil {
		return err
	}

	wrappedPty, err := wrappedParty.WrapParty(retrievePartyResponse.Party)
	if err != nil {
		return err
	}

	response.Party = *wrappedPty

	return nil
}

package jsonRpc

import (
	"gitlab.com/iotTracker/brain/log"
	"gitlab.com/iotTracker/brain/party"
	partyAdministrator "gitlab.com/iotTracker/brain/party/administrator"
	"gitlab.com/iotTracker/brain/search/wrappedIdentifier"
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

	getMyPartyResponse := partyAdministrator.GetMyPartyResponse{}
	if err := a.partyAdministrator.GetMyParty(&partyAdministrator.GetMyPartyRequest{
		Claims: claims,
	}, &getMyPartyResponse); err != nil {
		return err
	}

	response.Party = getMyPartyResponse.Party
	response.PartyType = getMyPartyResponse.PartyType

	return nil
}

type RetrievePartyRequest struct {
	PartyType  party.Type
	Identifier wrappedIdentifier.WrappedIdentifier
}

type RetrievePartyResponse struct {
	Party     party.Party
	PartyType party.Type
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

	retrievePartyResponse := partyAdministrator.RetrievePartyResponse{}
	if err := a.partyAdministrator.RetrieveParty(&partyAdministrator.RetrievePartyRequest{
		Claims:     claims,
		PartyType:  request.PartyType,
		Identifier: partyIdentifier,
	}, &retrievePartyResponse); err != nil {
		return err
	}

	response.Party = retrievePartyResponse.Party

	return nil
}

package jsonRpc

import (
	partyHandler "gitlab.com/iotTracker/brain/party/handler"
	"net/http"
	"gitlab.com/iotTracker/brain/security/wrappedClaims"
	"gitlab.com/iotTracker/brain/log"
	"gitlab.com/iotTracker/brain/party"
)

type adaptor struct {
	partyHandler partyHandler.Handler
}

func New(
	partyHandler partyHandler.Handler,
) *adaptor {
	return &adaptor{
		partyHandler: partyHandler,
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

	getMyPartyResponse := partyHandler.GetMyPartyResponse{}
	if err := a.partyHandler.GetMyParty(&partyHandler.GetMyPartyRequest{
		Claims: claims,
	}, &getMyPartyResponse); err != nil {
		return err
	}

	response.Party = getMyPartyResponse.Party
	response.PartyType = getMyPartyResponse.PartyType

	return nil
}

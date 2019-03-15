package jsonRpc

import (
	"gitlab.com/iotTracker/brain/log"
	"gitlab.com/iotTracker/brain/party"
	partyHandler "gitlab.com/iotTracker/brain/party/handler"
	"gitlab.com/iotTracker/brain/party/user"
	"gitlab.com/iotTracker/brain/security/wrappedClaims"
	"net/http"
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

type GetMyUserRequest struct{}

type GetMyUserResponse struct {
	User user.User `json:"user"`
}

func (a *adaptor) GetMyUser(r *http.Request, request *GetMyUserRequest, response *GetMyUserResponse) error {
	claims, err := wrappedClaims.UnwrapClaimsFromContext(r)
	if err != nil {
		log.Warn(err.Error())
		return err
	}

	getMyUserResponse := partyHandler.GetMyUserResponse{}
	if err := a.partyHandler.GetMyUser(&partyHandler.GetMyUserRequest{
		Claims: claims,
	}, &getMyUserResponse); err != nil {
		return err
	}

	response.User = getMyUserResponse.User

	return nil
}

package jsonRpc

import (
	"github.com/iot-my-world/brain/log"
	"github.com/iot-my-world/brain/party"
	partyAdministrator "github.com/iot-my-world/brain/party/administrator"
	"github.com/iot-my-world/brain/party/client"
	"github.com/iot-my-world/brain/party/company"
	wrappedParty "github.com/iot-my-world/brain/party/wrapped"
	wrappedIdentifier "github.com/iot-my-world/brain/search/identifier/wrapped"
	wrappedClaims "github.com/iot-my-world/brain/security/claims/wrapped"
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
	PartyType         party.Type
	WrappedIdentifier wrappedIdentifier.Wrapped
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

	retrievePartyResponse, err := a.partyAdministrator.RetrieveParty(&partyAdministrator.RetrievePartyRequest{
		Claims:     claims,
		PartyType:  request.PartyType,
		Identifier: request.WrappedIdentifier.Identifier,
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

type CreateAndInviteCompanyRequest struct {
	Company company.Company `json:"company"`
}

type CreateAndInviteCompanyResponse struct {
	RegistrationURLToken string `json:"registrationURLToken"`
}

func (a *adaptor) CreateAndInviteCompany(r *http.Request, request *CreateAndInviteCompanyRequest, response *CreateAndInviteCompanyResponse) error {
	createAndInviteCompanyResponse, err := a.partyAdministrator.CreateAndInviteCompany(&partyAdministrator.CreateAndInviteCompanyRequest{
		Company: request.Company,
	})
	if err != nil {
		return err
	}

	response.RegistrationURLToken = createAndInviteCompanyResponse.RegistrationURLToken

	return nil
}

type CreateAndInviteClientRequest struct {
	Client client.Client `json:"client"`
}

type CreateAndInviteResponse struct {
	RegistrationURLToken string `json:"registrationURLToken"`
}

func (a *adaptor) CreateAndInviteClient(r *http.Request, request *CreateAndInviteClientRequest, response *CreateAndInviteResponse) error {
	createAndInviteCompanyClientResponse, err := a.partyAdministrator.CreateAndInviteClient(&partyAdministrator.CreateAndInviteClientRequest{
		Client: request.Client,
	})
	if err != nil {
		return err
	}

	response.RegistrationURLToken = createAndInviteCompanyClientResponse.RegistrationURLToken

	return nil
}

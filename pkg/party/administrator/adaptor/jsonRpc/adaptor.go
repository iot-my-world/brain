package jsonRpc

import (
	"github.com/iot-my-world/brain/internal/log"
	"github.com/iot-my-world/brain/pkg/party"
	"github.com/iot-my-world/brain/pkg/party/administrator"
	client2 "github.com/iot-my-world/brain/pkg/party/client"
	company2 "github.com/iot-my-world/brain/pkg/party/company"
	"github.com/iot-my-world/brain/pkg/party/wrapped"
	wrappedIdentifier "github.com/iot-my-world/brain/pkg/search/identifier/wrapped"
	wrappedClaims "github.com/iot-my-world/brain/security/claims/wrapped"
	"net/http"
)

type adaptor struct {
	partyAdministrator administrator.Administrator
}

func New(
	partyAdministrator administrator.Administrator,
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

	getMyPartyResponse, err := a.partyAdministrator.GetMyParty(&administrator.GetMyPartyRequest{
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
	Party wrapped.Wrapped `json:"Party"`
}

func (a *adaptor) RetrieveParty(r *http.Request, request *RetrievePartyRequest, response *RetrievePartyResponse) error {
	claims, err := wrappedClaims.UnwrapClaimsFromContext(r)
	if err != nil {
		log.Warn(err.Error())
		return err
	}

	retrievePartyResponse, err := a.partyAdministrator.RetrieveParty(&administrator.RetrievePartyRequest{
		Claims:     claims,
		PartyType:  request.PartyType,
		Identifier: request.WrappedIdentifier.Identifier,
	})
	if err != nil {
		return err
	}

	wrappedPty, err := wrapped.WrapParty(retrievePartyResponse.Party)
	if err != nil {
		return err
	}

	response.Party = *wrappedPty

	return nil
}

type CreateAndInviteCompanyRequest struct {
	Company company2.Company `json:"company"`
}

type CreateAndInviteCompanyResponse struct {
	RegistrationURLToken string `json:"registrationURLToken"`
}

func (a *adaptor) CreateAndInviteCompany(r *http.Request, request *CreateAndInviteCompanyRequest, response *CreateAndInviteCompanyResponse) error {
	createAndInviteCompanyResponse, err := a.partyAdministrator.CreateAndInviteCompany(&administrator.CreateAndInviteCompanyRequest{
		Company: request.Company,
	})
	if err != nil {
		return err
	}

	response.RegistrationURLToken = createAndInviteCompanyResponse.RegistrationURLToken

	return nil
}

type CreateAndInviteClientRequest struct {
	Client client2.Client `json:"client"`
}

type CreateAndInviteResponse struct {
	RegistrationURLToken string `json:"registrationURLToken"`
}

func (a *adaptor) CreateAndInviteClient(r *http.Request, request *CreateAndInviteClientRequest, response *CreateAndInviteResponse) error {
	createAndInviteCompanyClientResponse, err := a.partyAdministrator.CreateAndInviteClient(&administrator.CreateAndInviteClientRequest{
		Client: request.Client,
	})
	if err != nil {
		return err
	}

	response.RegistrationURLToken = createAndInviteCompanyClientResponse.RegistrationURLToken

	return nil
}

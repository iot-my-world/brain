package jsonRpc

import (
	"errors"
	"github.com/iot-my-world/brain/internal/log"
	jsonRpcServiceProvider "github.com/iot-my-world/brain/pkg/api/jsonRpc/service/provider"
	partyRegistrar "github.com/iot-my-world/brain/pkg/party/registrar"
	"github.com/iot-my-world/brain/pkg/search/identifier/party"
	wrappedIdentifier "github.com/iot-my-world/brain/pkg/search/identifier/wrapped"
	wrappedClaims "github.com/iot-my-world/brain/pkg/security/claims/wrapped"
	humanUser "github.com/iot-my-world/brain/pkg/user/human"
	"net/http"
)

type adaptor struct {
	registrar partyRegistrar.Registrar
}

func New(
	registrar partyRegistrar.Registrar,
) *adaptor {
	return &adaptor{
		registrar: registrar,
	}
}

func (a *adaptor) Name() jsonRpcServiceProvider.Name {
	return jsonRpcServiceProvider.Name(partyRegistrar.ServiceProvider)
}

func (a *adaptor) MethodRequiresAuthorization(string) bool {
	return true
}

type InviteCompanyAdminUserRequest struct {
	WrappedCompanyIdentifier wrappedIdentifier.Wrapped `json:"companyIdentifier"`
}

type InviteCompanyAdminUserResponse struct {
	URLToken string `json:"urlToken"`
}

func (a *adaptor) InviteCompanyAdminUser(r *http.Request, request *InviteCompanyAdminUserRequest, response *InviteCompanyAdminUserResponse) error {
	claims, err := wrappedClaims.UnwrapClaimsFromContext(r)
	if err != nil {
		log.Error(err.Error())
		return err
	}

	inviteCompanyAdminUserResponse, err := a.registrar.InviteCompanyAdminUser(&partyRegistrar.InviteCompanyAdminUserRequest{
		Claims:            claims,
		CompanyIdentifier: request.WrappedCompanyIdentifier.Identifier,
	})
	if err != nil {
		return err
	}
	response.URLToken = inviteCompanyAdminUserResponse.URLToken
	return nil
}

type RegisterCompanyAdminUserRequest struct {
	User humanUser.User `json:"user"`
}

type RegisterCompanyAdminUserResponse struct {
	User humanUser.User `json:"user"`
}

func (a *adaptor) RegisterCompanyAdminUser(r *http.Request, request *RegisterCompanyAdminUserRequest, response *RegisterCompanyAdminUserResponse) error {
	claims, err := wrappedClaims.UnwrapClaimsFromContext(r)
	if err != nil {
		log.Error(err.Error())
		return err
	}

	registerUserResponse, err := a.registrar.RegisterCompanyAdminUser(&partyRegistrar.RegisterCompanyAdminUserRequest{
		Claims: claims,
		User:   request.User,
	})
	if err != nil {
		return err
	}

	response.User = registerUserResponse.User

	return nil
}

type RegisterCompanyUserRequest struct {
	User humanUser.User `json:"user"`
}

type RegisterCompanyUserResponse struct {
	User humanUser.User `json:"user"`
}

func (a *adaptor) RegisterCompanyUser(r *http.Request, request *RegisterCompanyUserRequest, response *RegisterCompanyUserResponse) error {
	claims, err := wrappedClaims.UnwrapClaimsFromContext(r)
	if err != nil {
		log.Error(err.Error())
		return err
	}

	registerResponse, err := a.registrar.RegisterCompanyUser(&partyRegistrar.RegisterCompanyUserRequest{
		Claims: claims,
		User:   request.User,
	})
	if err != nil {
		return err
	}

	response.User = registerResponse.User

	return nil
}

type InviteClientAdminUserRequest struct {
	WrappedClientIdentifier wrappedIdentifier.Wrapped `json:"clientIdentifier"`
}

type InviteClientAdminUserResponse struct {
	URLToken string `json:"urlToken"`
}

func (a *adaptor) InviteClientAdminUser(r *http.Request, request *InviteClientAdminUserRequest, response *InviteClientAdminUserResponse) error {

	claims, err := wrappedClaims.UnwrapClaimsFromContext(r)
	if err != nil {
		log.Error(err.Error())
		return err
	}

	inviteClientAdminUserResponse, err := a.registrar.InviteClientAdminUser(&partyRegistrar.InviteClientAdminUserRequest{
		Claims:           claims,
		ClientIdentifier: request.WrappedClientIdentifier.Identifier,
	})
	if err != nil {
		return err
	}
	response.URLToken = inviteClientAdminUserResponse.URLToken
	return nil
}

type RegisterClientAdminUserRequest struct {
	User humanUser.User `json:"user"`
}

type RegisterClientAdminUserResponse struct {
	User humanUser.User `json:"user"`
}

func (a *adaptor) RegisterClientAdminUser(r *http.Request, request *RegisterClientAdminUserRequest, response *RegisterClientAdminUserResponse) error {
	claims, err := wrappedClaims.UnwrapClaimsFromContext(r)
	if err != nil {
		log.Error(err.Error())
		return err
	}

	registerResponse, err := a.registrar.RegisterClientAdminUser(&partyRegistrar.RegisterClientAdminUserRequest{
		Claims: claims,
		User:   request.User,
	})
	if err != nil {
		return err
	}

	response.User = registerResponse.User

	return nil
}

type RegisterClientUserRequest struct {
	User humanUser.User `json:"user"`
}

type RegisterClientUserResponse struct {
	User humanUser.User `json:"user"`
}

func (a *adaptor) RegisterClientUser(r *http.Request, request *RegisterClientUserRequest, response *RegisterClientUserResponse) error {
	claims, err := wrappedClaims.UnwrapClaimsFromContext(r)
	if err != nil {
		log.Error(err.Error())
		return err
	}

	registerResponse, err := a.registrar.RegisterClientUser(&partyRegistrar.RegisterClientUserRequest{
		Claims: claims,
		User:   request.User,
	})
	if err != nil {
		return err
	}

	response.User = registerResponse.User

	return nil
}

type AreAdminsRegisteredRequest struct {
	WrappedPartyIdentifiers []wrappedIdentifier.Wrapped `json:"partyIdentifiers"`
}

type AreAdminsRegisteredResponse struct {
	Result map[string]bool `json:"result"`
}

func (a *adaptor) AreAdminsRegistered(r *http.Request, request *AreAdminsRegisteredRequest, response *AreAdminsRegisteredResponse) error {
	claims, err := wrappedClaims.UnwrapClaimsFromContext(r)
	if err != nil {
		log.Error(err.Error())
		return err
	}

	partyIdentifiers := make([]party.Identifier, 0)

	for i := range request.WrappedPartyIdentifiers {
		partyIdentifier, ok := request.WrappedPartyIdentifiers[i].Identifier.(party.Identifier)
		if !ok {
			return errors.New("could not cast identifier.Identifier to party.Identifier")
		}
		partyIdentifiers = append(partyIdentifiers, partyIdentifier)
	}

	areAdminsRegisteredResponse, err := a.registrar.AreAdminsRegistered(&partyRegistrar.AreAdminsRegisteredRequest{
		Claims:           claims,
		PartyIdentifiers: partyIdentifiers,
	})
	if err != nil {
		return err
	}

	response.Result = areAdminsRegisteredResponse.Result

	return nil
}

type InviteUserRequest struct {
	WrappedUserIdentifier wrappedIdentifier.Wrapped `json:"userIdentifier"`
}

type InviteUserResponse struct {
	URLToken string `json:"urlToken"`
}

func (a *adaptor) InviteUser(r *http.Request, request *InviteUserRequest, response *InviteUserResponse) error {
	claims, err := wrappedClaims.UnwrapClaimsFromContext(r)
	if err != nil {
		log.Error(err.Error())
		return err
	}

	userInviteResponse, err := a.registrar.InviteUser(&partyRegistrar.InviteUserRequest{
		Claims:         claims,
		UserIdentifier: request.WrappedUserIdentifier.Identifier,
	})
	if err != nil {
		return err
	}

	response.URLToken = userInviteResponse.URLToken

	return nil
}

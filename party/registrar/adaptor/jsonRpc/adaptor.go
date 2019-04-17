package jsonRpc

import (
	"errors"
	"gitlab.com/iotTracker/brain/log"
	"gitlab.com/iotTracker/brain/party/registrar"
	"gitlab.com/iotTracker/brain/search/identifier/party"
	wrappedIdentifier "gitlab.com/iotTracker/brain/search/identifier/wrapped"
	wrappedClaims "gitlab.com/iotTracker/brain/security/claims/wrapped"
	"gitlab.com/iotTracker/brain/user"
	"net/http"
)

type adaptor struct {
	registrar registrar.Registrar
}

func New(
	registrar registrar.Registrar,
) *adaptor {
	return &adaptor{
		registrar: registrar,
	}
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
		log.Warn(err.Error())
		return err
	}

	inviteCompanyAdminUserResponse, err := a.registrar.InviteCompanyAdminUser(&registrar.InviteCompanyAdminUserRequest{
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
	User user.User `json:"user"`
}

type RegisterCompanyAdminUserResponse struct {
	User user.User `json:"user"`
}

func (a *adaptor) RegisterCompanyAdminUser(r *http.Request, request *RegisterCompanyAdminUserRequest, response *RegisterCompanyAdminUserResponse) error {
	claims, err := wrappedClaims.UnwrapClaimsFromContext(r)
	if err != nil {
		log.Warn(err.Error())
		return err
	}

	registerCompanyAdminUserResponse, err := a.registrar.RegisterCompanyAdminUser(&registrar.RegisterCompanyAdminUserRequest{
		Claims: claims,
		User:   request.User,
	})
	if err != nil {
		log.Warn(err.Error())
		return err
	}

	response.User = registerCompanyAdminUserResponse.User

	return nil
}

type RegisterCompanyUserRequest struct {
	User user.User `json:"user"`
}

type RegisterCompanyUserResponse struct {
	User user.User `json:"user"`
}

func (a *adaptor) RegisterCompanyUser(r *http.Request, request *RegisterCompanyUserRequest, response *RegisterCompanyUserResponse) error {
	claims, err := wrappedClaims.UnwrapClaimsFromContext(r)
	if err != nil {
		log.Warn(err.Error())
		return err
	}

	registerCompanyUserResponse, err := a.registrar.RegisterCompanyUser(&registrar.RegisterCompanyUserRequest{
		Claims: claims,
		User:   request.User,
	})
	if err != nil {
		log.Warn(err.Error())
		return err
	}

	response.User = registerCompanyUserResponse.User

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
		log.Warn(err.Error())
		return err
	}

	inviteClientAdminUserResponse, err := a.registrar.InviteClientAdminUser(&registrar.InviteClientAdminUserRequest{
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
	User user.User `json:"user"`
}

type RegisterClientAdminUserResponse struct {
	User user.User `json:"user"`
}

func (a *adaptor) RegisterClientAdminUser(r *http.Request, request *RegisterClientAdminUserRequest, response *RegisterClientAdminUserResponse) error {
	claims, err := wrappedClaims.UnwrapClaimsFromContext(r)
	if err != nil {
		log.Warn(err.Error())
		return err
	}

	registerClientAdminUserResponse, err := a.registrar.RegisterClientAdminUser(&registrar.RegisterClientAdminUserRequest{
		Claims: claims,
		User:   request.User,
	})
	if err != nil {
		log.Warn(err.Error())
		return err
	}

	response.User = registerClientAdminUserResponse.User

	return nil
}

type RegisterClientUserRequest struct {
	User user.User `json:"user"`
}

type RegisterClientUserResponse struct {
	User user.User `json:"user"`
}

func (a *adaptor) RegisterClientUser(r *http.Request, request *RegisterClientUserRequest, response *RegisterClientUserResponse) error {
	claims, err := wrappedClaims.UnwrapClaimsFromContext(r)
	if err != nil {
		log.Warn(err.Error())
		return err
	}

	registerClientUserResponse, err := a.registrar.RegisterClientUser(&registrar.RegisterClientUserRequest{
		Claims: claims,
		User:   request.User,
	})
	if err != nil {
		log.Warn(err.Error())
		return err
	}

	response.User = registerClientUserResponse.User

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
		log.Warn(err.Error())
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

	areAdminsRegisteredResponse, err := a.registrar.AreAdminsRegistered(&registrar.AreAdminsRegisteredRequest{
		Claims:           claims,
		PartyIdentifiers: partyIdentifiers,
	})
	if err != nil {
		log.Warn(err.Error())
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
		log.Warn(err.Error())
		return err
	}

	userInviteResponse, err := a.registrar.InviteUser(&registrar.InviteUserRequest{
		Claims:         claims,
		UserIdentifier: request.WrappedUserIdentifier.Identifier,
	})
	if err != nil {
		log.Warn(err.Error())
		return err
	}

	response.URLToken = userInviteResponse.URLToken

	return nil
}

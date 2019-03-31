package jsonRpc

import (
	"gitlab.com/iotTracker/brain/log"
	partyIdentifier "gitlab.com/iotTracker/brain/search/identifier/party"
	"gitlab.com/iotTracker/brain/party/registrar"
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
	CompanyIdentifier wrappedIdentifier.Wrapped `json:"companyIdentifier"`
}

type InviteCompanyAdminUserResponse struct {
	URLToken string `json:"urlToken"`
}

func (a *adaptor) InviteCompanyAdminUser(r *http.Request, request *InviteCompanyAdminUserRequest, response *InviteCompanyAdminUserResponse) error {
	companyIdentifier, err := request.CompanyIdentifier.UnWrap()
	if err != nil {
		return err
	}

	claims, err := wrappedClaims.UnwrapClaimsFromContext(r)
	if err != nil {
		log.Warn(err.Error())
		return err
	}

	inviteCompanyAdminUserResponse, err := a.registrar.InviteCompanyAdminUser(&registrar.InviteCompanyAdminUserRequest{
		Claims:            claims,
		CompanyIdentifier: companyIdentifier,
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
	ClientIdentifier wrappedIdentifier.Wrapped `json:"clientIdentifier"`
}

type InviteClientAdminUserResponse struct {
	URLToken string `json:"urlToken"`
}

func (a *adaptor) InviteClientAdminUser(r *http.Request, request *InviteClientAdminUserRequest, response *InviteClientAdminUserResponse) error {
	clientIdentifier, err := request.ClientIdentifier.UnWrap()
	if err != nil {
		return err
	}

	claims, err := wrappedClaims.UnwrapClaimsFromContext(r)
	if err != nil {
		log.Warn(err.Error())
		return err
	}

	inviteClientAdminUserResponse, err := a.registrar.InviteClientAdminUser(&registrar.InviteClientAdminUserRequest{
		Claims:           claims,
		ClientIdentifier: clientIdentifier,
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
	PartyIdentifiers []partyIdentifier.Identifier `json:"partyIdentifiers"`
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

	areAdminsRegisteredResponse, err := a.registrar.AreAdminsRegistered(&registrar.AreAdminsRegisteredRequest{
		Claims:           claims,
		PartyIdentifiers: request.PartyIdentifiers,
	})
	if err != nil {
		log.Warn(err.Error())
		return err
	}

	response.Result = areAdminsRegisteredResponse.Result

	return nil
}

type InviteUserRequest struct {
	UserIdentifier wrappedIdentifier.Wrapped `json:"userIdentifier"`
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

	userId, err := request.UserIdentifier.UnWrap()
	if err != nil {
		return err
	}

	userInviteResponse, err := a.registrar.InviteUser(&registrar.InviteUserRequest{
		Claims:         claims,
		UserIdentifier: userId,
	})
	if err != nil {
		log.Warn(err.Error())
		return err
	}

	response.URLToken = userInviteResponse.URLToken

	return nil
}

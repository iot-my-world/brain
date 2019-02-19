package jsonRpc

import (
	"gitlab.com/iotTracker/brain/log"
	"gitlab.com/iotTracker/brain/party/registrar"
	"gitlab.com/iotTracker/brain/party/user"
	"gitlab.com/iotTracker/brain/search/wrappedIdentifier"
	"gitlab.com/iotTracker/brain/security/wrappedClaims"
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
	PartyIdentifier wrappedIdentifier.WrappedIdentifier `json:"partyIdentifier"`
}

type InviteCompanyAdminUserResponse struct {
}

func (a *adaptor) InviteCompanyAdminUser(r *http.Request, request *InviteCompanyAdminUserRequest, response *InviteCompanyAdminUserResponse) error {

	id, err := request.PartyIdentifier.UnWrap()
	if err != nil {
		return err
	}

	inviteCompanyAdminUserResponse := registrar.InviteCompanyAdminUserResponse{}
	if err := a.registrar.InviteCompanyAdminUser(&registrar.InviteCompanyAdminUserRequest{
		PartyIdentifier: id,
	},
		&inviteCompanyAdminUserResponse); err != nil {
		return err
	}

	return nil
}

type RegisterCompanyAdminUserRequest struct {
	User     user.User `json:"user"`
	Password string    `json:"password"`
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

	registerCompanyAdminUserResponse := registrar.RegisterCompanyAdminUserResponse{}
	if err := a.registrar.RegisterCompanyAdminUser(&registrar.RegisterCompanyAdminUserRequest{
		Claims:   claims,
		User:     request.User,
		Password: request.Password,
	},
		&registerCompanyAdminUserResponse); err != nil {
		log.Warn(err.Error())
		return err
	}

	response.User = registerCompanyAdminUserResponse.User

	return nil
}

type InviteClientAdminUserRequest struct {
	PartyIdentifier wrappedIdentifier.WrappedIdentifier `json:"partyIdentifier"`
}

type InviteClientAdminUserResponse struct {
}

func (a *adaptor) InviteClientAdminUser(r *http.Request, request *InviteClientAdminUserRequest, response *InviteClientAdminUserResponse) error {

	id, err := request.PartyIdentifier.UnWrap()
	if err != nil {
		return err
	}

	inviteClientAdminUserResponse := registrar.InviteClientAdminUserResponse{}
	if err := a.registrar.InviteClientAdminUser(&registrar.InviteClientAdminUserRequest{
		PartyIdentifier: id,
	},
		&inviteClientAdminUserResponse); err != nil {
		return err
	}

	return nil
}

type RegisterClientAdminUserRequest struct {
	User     user.User `json:"user"`
	Password string    `json:"password"`
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

	registerClientAdminUserResponse := registrar.RegisterClientAdminUserResponse{}
	if err := a.registrar.RegisterClientAdminUser(&registrar.RegisterClientAdminUserRequest{
		Claims:   claims,
		User:     request.User,
		Password: request.Password,
	},
		&registerClientAdminUserResponse); err != nil {
		log.Warn(err.Error())
		return err
	}

	response.User = registerClientAdminUserResponse.User

	return nil
}

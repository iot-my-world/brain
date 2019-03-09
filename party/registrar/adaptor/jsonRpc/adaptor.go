package jsonRpc

import (
	"gitlab.com/iotTracker/brain/log"
	"gitlab.com/iotTracker/brain/party/registrar"
	"gitlab.com/iotTracker/brain/party/user"
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
	User user.User `json:"user"`
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

	inviteCompanyAdminUserResponse := registrar.InviteCompanyAdminUserResponse{}
	if err := a.registrar.InviteCompanyAdminUser(&registrar.InviteCompanyAdminUserRequest{
		Claims: claims,
		User:   request.User,
	},
		&inviteCompanyAdminUserResponse); err != nil {
		return err
	}
	response.URLToken = inviteCompanyAdminUserResponse.URLToken
	return nil
}

type RegisterCompanyAdminUserRequest struct {
	User     user.User `json:"user"`
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
	},
		&registerCompanyAdminUserResponse); err != nil {
		log.Warn(err.Error())
		return err
	}

	response.User = registerCompanyAdminUserResponse.User

	return nil
}

type InviteCompanyUserRequest struct {
	User user.User `json:"user"`
}

type InviteCompanyUserResponse struct {
	URLToken string `json:"urlToken"`
}

func (a *adaptor) InviteCompanyUser(r *http.Request, request *InviteCompanyUserRequest, response *InviteCompanyUserResponse) error {
	claims, err := wrappedClaims.UnwrapClaimsFromContext(r)
	if err != nil {
		log.Warn(err.Error())
		return err
	}

	inviteCompanyUserResponse := registrar.InviteCompanyUserResponse{}
	if err := a.registrar.InviteCompanyUser(&registrar.InviteCompanyUserRequest{
		Claims: claims,
		User:   request.User,
	},
		&inviteCompanyUserResponse); err != nil {
		return err
	}
	response.URLToken = inviteCompanyUserResponse.URLToken
	return nil
}

type RegisterCompanyUserRequest struct {
	User     user.User `json:"user"`
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

	registerCompanyUserResponse := registrar.RegisterCompanyUserResponse{}
	if err := a.registrar.RegisterCompanyUser(&registrar.RegisterCompanyUserRequest{
		Claims:   claims,
		User:     request.User,
	},
		&registerCompanyUserResponse); err != nil {
		log.Warn(err.Error())
		return err
	}

	response.User = registerCompanyUserResponse.User

	return nil
}

type InviteClientAdminUserRequest struct {
	User user.User `json:"user"`
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

	inviteClientAdminUserResponse := registrar.InviteClientAdminUserResponse{}
	if err := a.registrar.InviteClientAdminUser(&registrar.InviteClientAdminUserRequest{
		Claims: claims,
		User:   request.User,
	},
		&inviteClientAdminUserResponse); err != nil {
		return err
	}
	response.URLToken = inviteClientAdminUserResponse.URLToken
	return nil
}

type RegisterClientAdminUserRequest struct {
	User     user.User `json:"user"`
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
	},
		&registerClientAdminUserResponse); err != nil {
		log.Warn(err.Error())
		return err
	}

	response.User = registerClientAdminUserResponse.User

	return nil
}

type InviteClientUserRequest struct {
	User user.User `json:"user"`
}

type InviteClientUserResponse struct {
	URLToken string `json:"urlToken"`
}

func (a *adaptor) InviteClientUser(r *http.Request, request *InviteClientUserRequest, response *InviteClientUserResponse) error {
	claims, err := wrappedClaims.UnwrapClaimsFromContext(r)
	if err != nil {
		log.Warn(err.Error())
		return err
	}

	inviteClientUserResponse := registrar.InviteClientUserResponse{}
	if err := a.registrar.InviteClientUser(&registrar.InviteClientUserRequest{
		Claims: claims,
		User:   request.User,
	},
		&inviteClientUserResponse); err != nil {
		return err
	}
	response.URLToken = inviteClientUserResponse.URLToken
	return nil
}

type RegisterClientUserRequest struct {
	User     user.User `json:"user"`
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

	registerClientUserResponse := registrar.RegisterClientUserResponse{}
	if err := a.registrar.RegisterClientUser(&registrar.RegisterClientUserRequest{
		Claims:   claims,
		User:     request.User,
	},
		&registerClientUserResponse); err != nil {
		log.Warn(err.Error())
		return err
	}

	response.User = registerClientUserResponse.User

	return nil
}

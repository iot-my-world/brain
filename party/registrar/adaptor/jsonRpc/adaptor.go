package jsonRpc

import (
	"gitlab.com/iotTracker/brain/party/registrar"
	"gitlab.com/iotTracker/brain/search/wrappedIdentifier"
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

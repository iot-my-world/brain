package jsonRpc

import (
	tk102Administrator "gitlab.com/iotTracker/brain/tracker/device/tk102/administrator"
	"net/http"
	"gitlab.com/iotTracker/brain/search/wrappedIdentifier"
	"gitlab.com/iotTracker/brain/party"
	"gitlab.com/iotTracker/brain/search/identifier"
	"gitlab.com/iotTracker/brain/security/wrappedClaims"
	"gitlab.com/iotTracker/brain/security/claims"
	"gitlab.com/iotTracker/brain/log"
)

type adaptor struct {
	administrator tk102Administrator.Administrator
}

func New(administrator tk102Administrator.Administrator) *adaptor {
	return &adaptor{
		administrator: administrator,
	}
}

type ChangeOwnerRequest struct {
	TK102Identifier    wrappedIdentifier.WrappedIdentifier `json:"tk102Identifier"`
	NewOwnerPartyType  party.Type                          `json:"newOwnerPartyType"`
	NewOwnerIdentifier wrappedIdentifier.WrappedIdentifier `json:"newOwnerIdentifier"`
}

type ChangeOwnerResponse struct {
}

func (a *adaptor) ChangeOwner(r *http.Request, request *ChangeOwnerRequest, response *ChangeOwnerResponse) error {
	var claims claims.Claims
	var tk102Id identifier.Identifier
	var newOwnerId identifier.Identifier
	var err error

	claims, err = wrappedClaims.UnwrapClaimsFromContext(r)
	if err != nil {
		log.Warn(err.Error())
		return err
	}

	tk102Id, err = request.TK102Identifier.UnWrap()
	if err != nil {
		return err
	}
	newOwnerId, err = request.NewOwnerIdentifier.UnWrap()
	if err != nil {
		return err
	}

	changeOwnerResponse := tk102Administrator.ChangeOwnerResponse{}
	if err := a.administrator.ChangeOwner(&tk102Administrator.ChangeOwnerRequest{
		Claims:             claims,
		Identifier:         tk102Id,
		NewOwnerPartyType:  request.NewOwnerPartyType,
		NewOwnerIdentifier: newOwnerId,
	},
		&changeOwnerResponse); err != nil {
		return err
	}

	return nil
}

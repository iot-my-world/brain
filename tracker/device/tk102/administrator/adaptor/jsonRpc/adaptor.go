package jsonRpc

import (
	"gitlab.com/iotTracker/brain/log"
	"gitlab.com/iotTracker/brain/party"
	"gitlab.com/iotTracker/brain/search/identifier"
	"gitlab.com/iotTracker/brain/search/wrappedIdentifier"
	"gitlab.com/iotTracker/brain/security/claims"
	"gitlab.com/iotTracker/brain/security/wrappedClaims"
	"gitlab.com/iotTracker/brain/tracker/device/tk102"
	tk102Administrator "gitlab.com/iotTracker/brain/tracker/device/tk102/administrator"
	"net/http"
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
	TK102 tk102.TK102 `json:"tk102"`
}

func (a *adaptor) ChangeOwner(r *http.Request, request *ChangeOwnerRequest, response *ChangeOwnerResponse) error {
	var loginClaims claims.Claims
	var tk102Id identifier.Identifier
	var newOwnerId identifier.Identifier
	var err error

	loginClaims, err = wrappedClaims.UnwrapClaimsFromContext(r)
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
		Claims:             loginClaims,
		TK02Identifier:     tk102Id,
		NewOwnerPartyType:  request.NewOwnerPartyType,
		NewOwnerIdentifier: newOwnerId,
	},
		&changeOwnerResponse); err != nil {
		return err
	}

	response.TK102 = changeOwnerResponse.TK102

	return nil
}

type ChangeAssignedRequest struct {
	TK102Identifier       wrappedIdentifier.WrappedIdentifier `json:"tk102Identifier"`
	NewAssignedPartyType  party.Type                          `json:"newAssignedPartyType"`
	NewAssignedIdentifier wrappedIdentifier.WrappedIdentifier `json:"newAssignedIdentifier"`
}

type ChangeAssignedResponse struct {
	TK102 tk102.TK102 `json:"tk102"`
}

func (a *adaptor) ChangeAssigned(r *http.Request, request *ChangeAssignedRequest, response *ChangeAssignedResponse) error {
	var loginClaims claims.Claims
	var tk102Id identifier.Identifier
	var newAssignedId identifier.Identifier
	var err error

	loginClaims, err = wrappedClaims.UnwrapClaimsFromContext(r)
	if err != nil {
		log.Warn(err.Error())
		return err
	}

	tk102Id, err = request.TK102Identifier.UnWrap()
	if err != nil {
		return err
	}
	newAssignedId, err = request.NewAssignedIdentifier.UnWrap()
	if err != nil {
		return err
	}

	changeAssignedResponse := tk102Administrator.ChangeAssignedResponse{}
	if err := a.administrator.ChangeAssigned(&tk102Administrator.ChangeAssignedRequest{
		Claims:                loginClaims,
		TK02Identifier:        tk102Id,
		NewAssignedPartyType:  request.NewAssignedPartyType,
		NewAssignedIdentifier: newAssignedId,
	},
		&changeAssignedResponse); err != nil {
		return err
	}

	response.TK102 = changeAssignedResponse.TK102

	return nil
}

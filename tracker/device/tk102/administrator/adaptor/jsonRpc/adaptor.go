package jsonRpc

import (
	"gitlab.com/iotTracker/brain/log"
	"gitlab.com/iotTracker/brain/security/wrappedClaims"
	"gitlab.com/iotTracker/brain/tracker/device/tk102"
	tk102Administrator "gitlab.com/iotTracker/brain/tracker/device/tk102/administrator"
	"net/http"
)

type Adaptor struct {
	administrator tk102Administrator.Administrator
}

func New(administrator tk102Administrator.Administrator) *Adaptor {
	return &Adaptor{
		administrator: administrator,
	}
}

type ChangeOwnershipAndAssignmentRequest struct {
	TK102 tk102.TK102 `json:"tk102"`
}

type ChangeOwnershipAndAssignmentResponse struct {
	TK102 tk102.TK102 `json:"tk102"`
}

func (a *Adaptor) ChangeOwnershipAndAssignment(r *http.Request, request *ChangeOwnershipAndAssignmentRequest, response *ChangeOwnershipAndAssignmentResponse) error {
	claims, err := wrappedClaims.UnwrapClaimsFromContext(r)
	if err != nil {
		log.Warn(err.Error())
		return err
	}

	changeOwnershipAndAssignmentResponse := tk102Administrator.ChangeOwnershipAndAssignmentResponse{}
	if err := a.administrator.ChangeOwnershipAndAssignment(&tk102Administrator.ChangeOwnershipAndAssignmentRequest{
		Claims: claims,
		TK102:  request.TK102,
	},
		&changeOwnershipAndAssignmentResponse); err != nil {
		return err
	}

	response.TK102 = changeOwnershipAndAssignmentResponse.TK102

	return nil
}

type CreateRequest struct {
	TK102 tk102.TK102 `json:"tk102"`
}

type CreateResponse struct {
	TK102 tk102.TK102 `json:"tk102"`
}

func (a *Adaptor) Create(r *http.Request, request *CreateRequest, response *CreateResponse) error {
	claims, err := wrappedClaims.UnwrapClaimsFromContext(r)
	if err != nil {
		log.Warn(err.Error())
		return err
	}

	return nil
}

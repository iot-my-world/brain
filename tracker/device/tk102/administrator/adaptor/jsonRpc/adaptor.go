package jsonRpc

import (
	"gitlab.com/iotTracker/brain/log"
	"gitlab.com/iotTracker/brain/security/wrappedClaims"
	"gitlab.com/iotTracker/brain/tracker/device/tk102"
	tk102Administrator "gitlab.com/iotTracker/brain/tracker/device/tk102/administrator"
	"net/http"
)

// Adaptor for tk102 device administrator for access via json rpc
type Adaptor struct {
	administrator tk102Administrator.Administrator
}

// New tk102 device administrator json rpc adaptor
func New(administrator tk102Administrator.Administrator) *Adaptor {
	return &Adaptor{
		administrator: administrator,
	}
}

// ChangeOwnershipAndAssignmentRequest contains parameters for a change ownership and assignment operation
type ChangeOwnershipAndAssignmentRequest struct {
	TK102 tk102.TK102 `json:"tk102"`
}

// ChangeOwnershipAndAssignmentResponse contains the tk102 device with changed ownership and assignment
type ChangeOwnershipAndAssignmentResponse struct {
	TK102 tk102.TK102 `json:"tk102"`
}

// ChangeOwnershipAndAssignment of a TK102 device using the tk102 device administrator
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

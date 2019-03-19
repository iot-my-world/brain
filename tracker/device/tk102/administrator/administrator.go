package administrator

import (
	"gitlab.com/iotTracker/brain/security/claims"
	"gitlab.com/iotTracker/brain/tracker/device/tk102"
)

type Administrator interface {
	ChangeOwnershipAndAssignment(request *ChangeOwnershipAndAssignmentRequest, response *ChangeOwnershipAndAssignmentResponse) error
	Create(request *CreateRequest, response *CreateResponse) error
}

type ChangeOwnershipAndAssignmentRequest struct {
	Claims claims.Claims
	TK102  tk102.TK102
}

type ChangeOwnershipAndAssignmentResponse struct {
	TK102 tk102.TK102
}

type CreateRequest struct {
	Claims claims.Claims
	TK102  tk102.TK102
}

type CreateResponse struct {
	TK102 tk102.TK102
}

package administrator

import (
	"gitlab.com/iotTracker/brain/security/claims"
	"gitlab.com/iotTracker/brain/tracker/device/tk102"
)

// Administrator performs administrative tasks on a TK102 Device
type Administrator interface {
	ChangeOwnershipAndAssignment(request *ChangeOwnershipAndAssignmentRequest, response *ChangeOwnershipAndAssignmentResponse) error
}

// ChangeOwnershipAndAssignmentRequest is the Administrator's ChangeOwnershipAndAssignment request object
type ChangeOwnershipAndAssignmentRequest struct {
	Claims claims.Claims
	TK102  tk102.TK102
}

// ChangeOwnershipAndAssignmentResponse is the Administrator's ChangeOwnershipAndAssignment response object
type ChangeOwnershipAndAssignmentResponse struct {
	TK102 tk102.TK102
}

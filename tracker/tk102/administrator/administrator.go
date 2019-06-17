package administrator

import (
	"github.com/iot-my-world/brain/security/claims"
	"github.com/iot-my-world/brain/tracker/tk102"
)

type Administrator interface {
	ChangeOwnershipAndAssignment(request *ChangeOwnershipAndAssignmentRequest) (*ChangeOwnershipAndAssignmentResponse, error)
	Create(request *CreateRequest) (*CreateResponse, error)
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

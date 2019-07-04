package administrator

import (
	tk1022 "github.com/iot-my-world/brain/pkg/tracker/tk102"
	"github.com/iot-my-world/brain/security/claims"
)

type Administrator interface {
	ChangeOwnershipAndAssignment(request *ChangeOwnershipAndAssignmentRequest) (*ChangeOwnershipAndAssignmentResponse, error)
	Create(request *CreateRequest) (*CreateResponse, error)
}

type ChangeOwnershipAndAssignmentRequest struct {
	Claims claims.Claims
	TK102  tk1022.TK102
}

type ChangeOwnershipAndAssignmentResponse struct {
	TK102 tk1022.TK102
}

type CreateRequest struct {
	Claims claims.Claims
	TK102  tk1022.TK102
}

type CreateResponse struct {
	TK102 tk1022.TK102
}

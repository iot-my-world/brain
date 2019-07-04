package jsonRpc

import (
	"github.com/iot-my-world/brain/internal/log"
	tk1022 "github.com/iot-my-world/brain/pkg/tracker/tk102"
	"github.com/iot-my-world/brain/pkg/tracker/tk102/administrator"
	wrappedClaims "github.com/iot-my-world/brain/security/claims/wrapped"
	"net/http"
)

type Adaptor struct {
	administrator administrator.Administrator
}

func New(administrator administrator.Administrator) *Adaptor {
	return &Adaptor{
		administrator: administrator,
	}
}

type ChangeOwnershipAndAssignmentRequest struct {
	TK102 tk1022.TK102 `json:"tk102"`
}

type ChangeOwnershipAndAssignmentResponse struct {
	TK102 tk1022.TK102 `json:"tk102"`
}

func (a *Adaptor) ChangeOwnershipAndAssignment(r *http.Request, request *ChangeOwnershipAndAssignmentRequest, response *ChangeOwnershipAndAssignmentResponse) error {
	claims, err := wrappedClaims.UnwrapClaimsFromContext(r)
	if err != nil {
		log.Warn(err.Error())
		return err
	}

	changeOwnershipAndAssignmentResponse, err := a.administrator.ChangeOwnershipAndAssignment(&administrator.ChangeOwnershipAndAssignmentRequest{
		Claims: claims,
		TK102:  request.TK102,
	})
	if err != nil {
		return err
	}

	response.TK102 = changeOwnershipAndAssignmentResponse.TK102

	return nil
}

type CreateRequest struct {
	TK102 tk1022.TK102 `json:"tk102"`
}

type CreateResponse struct {
	TK102 tk1022.TK102 `json:"tk102"`
}

func (a *Adaptor) Create(r *http.Request, request *CreateRequest, response *CreateResponse) error {
	claims, err := wrappedClaims.UnwrapClaimsFromContext(r)
	if err != nil {
		log.Warn(err.Error())
		return err
	}

	createResponse, err := a.administrator.Create(&administrator.CreateRequest{
		Claims: claims,
		TK102:  request.TK102,
	})
	if err != nil {
		return err
	}

	response.TK102 = createResponse.TK102

	return nil
}

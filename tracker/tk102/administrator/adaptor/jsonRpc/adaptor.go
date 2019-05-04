package jsonRpc

import (
	"gitlab.com/iotTracker/brain/log"
	wrappedClaims "gitlab.com/iotTracker/brain/security/claims/wrapped"
	"gitlab.com/iotTracker/brain/tracker/tk102"
	tk102DeviceAdministrator "gitlab.com/iotTracker/brain/tracker/tk102/administrator"
	"net/http"
)

type Adaptor struct {
	administrator tk102DeviceAdministrator.Administrator
}

func New(administrator tk102DeviceAdministrator.Administrator) *Adaptor {
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

	changeOwnershipAndAssignmentResponse, err := a.administrator.ChangeOwnershipAndAssignment(&tk102DeviceAdministrator.ChangeOwnershipAndAssignmentRequest{
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

	createResponse, err := a.administrator.Create(&tk102DeviceAdministrator.CreateRequest{
		Claims: claims,
		TK102:  request.TK102,
	})
	if err != nil {
		return err
	}

	response.TK102 = createResponse.TK102

	return nil
}
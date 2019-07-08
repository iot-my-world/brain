package jsonRpc

import (
	"github.com/iot-my-world/brain/internal/log"
	"github.com/iot-my-world/brain/pkg/device/sigbug"
	"github.com/iot-my-world/brain/pkg/device/sigbug/administrator"
	wrappedClaims "github.com/iot-my-world/brain/pkg/security/claims/wrapped"
	"net/http"
)

type adaptor struct {
	administrator administrator.Administrator
}

func New(administrator administrator.Administrator) *adaptor {
	return &adaptor{
		administrator: administrator,
	}
}

type CreateRequest struct {
	Sigbug sigbug.Sigbug `json:"sigbug"`
}

type CreateResponse struct {
	Sigbug sigbug.Sigbug `json:"sigbug"`
}

func (a *adaptor) Create(r *http.Request, request *CreateRequest, response *CreateResponse) error {
	claims, err := wrappedClaims.UnwrapClaimsFromContext(r)
	if err != nil {
		log.Warn(err.Error())
		return err
	}

	createResponse, err := a.administrator.Create(&administrator.CreateRequest{
		Claims: claims,
		Sigbug: request.Sigbug,
	})
	if err != nil {
		return err
	}

	response.Sigbug = createResponse.Sigbug

	return nil
}

type UpdateAllowedFieldsRequest struct {
	Sigbug sigbug.Sigbug `json:"sigbug"`
}

type UpdateAllowedFieldsResponse struct {
	Sigbug sigbug.Sigbug `json:"sigbug"`
}

func (a *adaptor) UpdateAllowedFields(r *http.Request, request *UpdateAllowedFieldsRequest, response *UpdateAllowedFieldsResponse) error {
	claims, err := wrappedClaims.UnwrapClaimsFromContext(r)
	if err != nil {
		log.Warn(err.Error())
		return err
	}

	updateAllowedFieldsResponse, err := a.administrator.UpdateAllowedFields(&administrator.UpdateAllowedFieldsRequest{
		Claims: claims,
		Sigbug: request.Sigbug,
	})
	if err != nil {
		return err
	}

	response.Sigbug = updateAllowedFieldsResponse.Sigbug

	return nil
}

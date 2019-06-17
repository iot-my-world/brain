package jsonRpc

import (
	"github.com/iot-my-world/brain/log"
	"github.com/iot-my-world/brain/party/system"
	systemAdministrator "github.com/iot-my-world/brain/party/system/administrator"
	wrappedClaims "github.com/iot-my-world/brain/security/claims/wrapped"
	"net/http"
)

type adaptor struct {
	systemAdministrator systemAdministrator.Administrator
}

func New(
	systemAdministrator systemAdministrator.Administrator,
) *adaptor {
	return &adaptor{
		systemAdministrator: systemAdministrator,
	}
}

type UpdateAllowedFieldsRequest struct {
	System system.System `json:"system"`
}

type UpdateAllowedFieldsResponse struct {
	System system.System `json:"system"`
}

func (a *adaptor) UpdateAllowedFields(r *http.Request, request *UpdateAllowedFieldsRequest, response *UpdateAllowedFieldsResponse) error {
	claims, err := wrappedClaims.UnwrapClaimsFromContext(r)
	if err != nil {
		log.Warn(err.Error())
		return err
	}

	updateAllowedFieldsResponse, err := a.systemAdministrator.UpdateAllowedFields(&systemAdministrator.UpdateAllowedFieldsRequest{
		Claims: claims,
		System: request.System,
	})
	if err != nil {
		return err
	}

	response.System = updateAllowedFieldsResponse.System

	return nil
}

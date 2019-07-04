package jsonRpc

import (
	"github.com/iot-my-world/brain/internal/log"
	system2 "github.com/iot-my-world/brain/pkg/party/system"
	"github.com/iot-my-world/brain/pkg/party/system/administrator"
	wrappedClaims "github.com/iot-my-world/brain/security/claims/wrapped"
	"net/http"
)

type adaptor struct {
	systemAdministrator administrator.Administrator
}

func New(
	systemAdministrator administrator.Administrator,
) *adaptor {
	return &adaptor{
		systemAdministrator: systemAdministrator,
	}
}

type UpdateAllowedFieldsRequest struct {
	System system2.System `json:"system"`
}

type UpdateAllowedFieldsResponse struct {
	System system2.System `json:"system"`
}

func (a *adaptor) UpdateAllowedFields(r *http.Request, request *UpdateAllowedFieldsRequest, response *UpdateAllowedFieldsResponse) error {
	claims, err := wrappedClaims.UnwrapClaimsFromContext(r)
	if err != nil {
		log.Warn(err.Error())
		return err
	}

	updateAllowedFieldsResponse, err := a.systemAdministrator.UpdateAllowedFields(&administrator.UpdateAllowedFieldsRequest{
		Claims: claims,
		System: request.System,
	})
	if err != nil {
		return err
	}

	response.System = updateAllowedFieldsResponse.System

	return nil
}

package jsonRpc

import (
	"github.com/iot-my-world/brain/log"
	wrappedClaims "github.com/iot-my-world/brain/security/claims/wrapped"
	"github.com/iot-my-world/brain/tracker/sf001"
	sf001DeviceAdministrator "github.com/iot-my-world/brain/tracker/sf001/administrator"
	"net/http"
)

type adaptor struct {
	administrator sf001DeviceAdministrator.Administrator
}

func New(administrator sf001DeviceAdministrator.Administrator) *adaptor {
	return &adaptor{
		administrator: administrator,
	}
}

type CreateRequest struct {
	SF001 sf001.SF001 `json:"sf001"`
}

type CreateResponse struct {
	SF001 sf001.SF001 `json:"sf001"`
}

func (a *adaptor) Create(r *http.Request, request *CreateRequest, response *CreateResponse) error {
	claims, err := wrappedClaims.UnwrapClaimsFromContext(r)
	if err != nil {
		log.Warn(err.Error())
		return err
	}

	createResponse, err := a.administrator.Create(&sf001DeviceAdministrator.CreateRequest{
		Claims: claims,
		SF001:  request.SF001,
	})
	if err != nil {
		return err
	}

	response.SF001 = createResponse.SF001

	return nil
}

type UpdateAllowedFieldsRequest struct {
	SF001 sf001.SF001 `json:"sf001"`
}

type UpdateAllowedFieldsResponse struct {
	SF001 sf001.SF001 `json:"sf001"`
}

func (a *adaptor) UpdateAllowedFields(r *http.Request, request *UpdateAllowedFieldsRequest, response *UpdateAllowedFieldsResponse) error {
	claims, err := wrappedClaims.UnwrapClaimsFromContext(r)
	if err != nil {
		log.Warn(err.Error())
		return err
	}

	updateAllowedFieldsResponse, err := a.administrator.UpdateAllowedFields(&sf001DeviceAdministrator.UpdateAllowedFieldsRequest{
		Claims: claims,
		SF001:  request.SF001,
	})
	if err != nil {
		return err
	}

	response.SF001 = updateAllowedFieldsResponse.SF001

	return nil
}

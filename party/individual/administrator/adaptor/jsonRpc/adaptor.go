package jsonRpc

import (
	"github.com/iot-my-world/brain/log"
	"github.com/iot-my-world/brain/party/individual"
	sf001DeviceAdministrator "github.com/iot-my-world/brain/party/individual/administrator"
	wrappedClaims "github.com/iot-my-world/brain/security/claims/wrapped"
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
	Individual individual.Individual `json:"individual"`
}

type CreateResponse struct {
	Individual individual.Individual `json:"individual"`
}

func (a *adaptor) Create(r *http.Request, request *CreateRequest, response *CreateResponse) error {
	claims, err := wrappedClaims.UnwrapClaimsFromContext(r)
	if err != nil {
		log.Warn(err.Error())
		return err
	}

	createResponse, err := a.administrator.Create(&sf001DeviceAdministrator.CreateRequest{
		Claims:     claims,
		Individual: request.Individual,
	})
	if err != nil {
		return err
	}

	response.Individual = createResponse.Individual

	return nil
}

type UpdateAllowedFieldsRequest struct {
	Individual individual.Individual `json:"individual"`
}

type UpdateAllowedFieldsResponse struct {
	Individual individual.Individual `json:"individual"`
}

func (a *adaptor) UpdateAllowedFields(r *http.Request, request *UpdateAllowedFieldsRequest, response *UpdateAllowedFieldsResponse) error {
	claims, err := wrappedClaims.UnwrapClaimsFromContext(r)
	if err != nil {
		log.Warn(err.Error())
		return err
	}

	updateAllowedFieldsResponse, err := a.administrator.UpdateAllowedFields(&sf001DeviceAdministrator.UpdateAllowedFieldsRequest{
		Claims:     claims,
		Individual: request.Individual,
	})
	if err != nil {
		return err
	}

	response.Individual = updateAllowedFieldsResponse.Individual

	return nil
}

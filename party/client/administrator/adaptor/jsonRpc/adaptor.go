package jsonRpc

import (
	"github.com/iot-my-world/brain/log"
	"github.com/iot-my-world/brain/party/client"
	clientAdministrator "github.com/iot-my-world/brain/party/client/administrator"
	wrappedClaims "github.com/iot-my-world/brain/security/claims/wrapped"
	"net/http"
)

type adaptor struct {
	clientAdministrator clientAdministrator.Administrator
}

func New(
	clientAdministrator clientAdministrator.Administrator,
) *adaptor {
	return &adaptor{
		clientAdministrator: clientAdministrator,
	}
}

type UpdateAllowedFieldsRequest struct {
	Client client.Client `json:"client"`
}

type UpdateAllowedFieldsResponse struct {
	Client client.Client `json:"client"`
}

func (a *adaptor) UpdateAllowedFields(r *http.Request, request *UpdateAllowedFieldsRequest, response *UpdateAllowedFieldsResponse) error {
	claims, err := wrappedClaims.UnwrapClaimsFromContext(r)
	if err != nil {
		log.Warn(err.Error())
		return err
	}

	updateAllowedFieldsResponse, err := a.clientAdministrator.UpdateAllowedFields(&clientAdministrator.UpdateAllowedFieldsRequest{
		Claims: claims,
		Client: request.Client,
	})
	if err != nil {
		return err
	}

	response.Client = updateAllowedFieldsResponse.Client

	return nil
}

type CreateRequest struct {
	Client client.Client `json:"client"`
}

type CreateResponse struct {
	Client client.Client `json:"client"`
}

func (a *adaptor) Create(r *http.Request, request *CreateRequest, response *CreateResponse) error {
	claims, err := wrappedClaims.UnwrapClaimsFromContext(r)
	if err != nil {
		log.Warn(err.Error())
		return err
	}

	createResponse, err := a.clientAdministrator.Create(&clientAdministrator.CreateRequest{
		Claims: claims,
		Client: request.Client,
	})
	if err != nil {
		return err
	}

	response.Client = createResponse.Client

	return nil
}

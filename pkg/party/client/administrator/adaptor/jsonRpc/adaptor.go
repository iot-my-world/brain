package jsonRpc

import (
	"github.com/iot-my-world/brain/internal/log"
	client2 "github.com/iot-my-world/brain/pkg/party/client"
	"github.com/iot-my-world/brain/pkg/party/client/administrator"
	wrappedIdentifier "github.com/iot-my-world/brain/pkg/search/identifier/wrapped"
	wrappedClaims "github.com/iot-my-world/brain/security/claims/wrapped"
	"net/http"
)

type adaptor struct {
	clientAdministrator administrator.Administrator
}

func New(
	clientAdministrator administrator.Administrator,
) *adaptor {
	return &adaptor{
		clientAdministrator: clientAdministrator,
	}
}

type UpdateAllowedFieldsRequest struct {
	Client client2.Client `json:"client"`
}

type UpdateAllowedFieldsResponse struct {
	Client client2.Client `json:"client"`
}

func (a *adaptor) UpdateAllowedFields(r *http.Request, request *UpdateAllowedFieldsRequest, response *UpdateAllowedFieldsResponse) error {
	claims, err := wrappedClaims.UnwrapClaimsFromContext(r)
	if err != nil {
		log.Warn(err.Error())
		return err
	}

	updateAllowedFieldsResponse, err := a.clientAdministrator.UpdateAllowedFields(&administrator.UpdateAllowedFieldsRequest{
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
	Client client2.Client `json:"client"`
}

type CreateResponse struct {
	Client client2.Client `json:"client"`
}

func (a *adaptor) Create(r *http.Request, request *CreateRequest, response *CreateResponse) error {
	claims, err := wrappedClaims.UnwrapClaimsFromContext(r)
	if err != nil {
		log.Warn(err.Error())
		return err
	}

	createResponse, err := a.clientAdministrator.Create(&administrator.CreateRequest{
		Claims: claims,
		Client: request.Client,
	})
	if err != nil {
		return err
	}

	response.Client = createResponse.Client

	return nil
}

type DeleteRequest struct {
	ClientIdentifier wrappedIdentifier.Wrapped `json:"clientIdentifier"`
}

type DeleteResponse struct {
}

func (a *adaptor) Delete(r *http.Request, request *DeleteRequest, response *DeleteResponse) error {
	claims, err := wrappedClaims.UnwrapClaimsFromContext(r)
	if err != nil {
		log.Warn(err.Error())
		return err
	}

	if _, err := a.clientAdministrator.Delete(&administrator.DeleteRequest{
		Claims:           claims,
		ClientIdentifier: request.ClientIdentifier.Identifier,
	}); err != nil {
		return err
	}

	return nil
}

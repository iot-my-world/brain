package jsonRpc

import (
	"gitlab.com/iotTracker/brain/log"
	"gitlab.com/iotTracker/brain/party/client"
	clientAdministrator "gitlab.com/iotTracker/brain/party/client/administrator"
	wrappedClaims "gitlab.com/iotTracker/brain/security/claims/wrapped"
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
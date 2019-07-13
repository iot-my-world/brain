package jsonRpc

import (
	"github.com/iot-my-world/brain/internal/log"
	jsonRpcServiceProvider "github.com/iot-my-world/brain/pkg/api/jsonRpc/service/provider"
	wrappedClaims "github.com/iot-my-world/brain/pkg/security/claims/wrapped"
	"github.com/iot-my-world/brain/pkg/sigfox/backend"
	"github.com/iot-my-world/brain/pkg/sigfox/backend/administrator"
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

func (a *adaptor) Name() jsonRpcServiceProvider.Name {
	return jsonRpcServiceProvider.Name(administrator.ServiceProvider)
}

func (a *adaptor) MethodRequiresAuthorization(string) bool {
	return true
}

type CreateRequest struct {
	Backend backend.Backend `json:"backend"`
}

type CreateResponse struct {
	Backend backend.Backend `json:"backend"`
}

func (a *adaptor) Create(r *http.Request, request *CreateRequest, response *CreateResponse) error {
	claims, err := wrappedClaims.UnwrapClaimsFromContext(r)
	if err != nil {
		log.Warn(err.Error())
		return err
	}

	createResponse, err := a.administrator.Create(&administrator.CreateRequest{
		Claims:  claims,
		Backend: request.Backend,
	})
	if err != nil {
		return err
	}

	response.Backend = createResponse.Backend

	return nil
}

type UpdateAllowedFieldsRequest struct {
	Backend backend.Backend `json:"backend"`
}

type UpdateAllowedFieldsResponse struct {
	Backend backend.Backend `json:"backend"`
}

func (a *adaptor) UpdateAllowedFields(r *http.Request, request *UpdateAllowedFieldsRequest, response *UpdateAllowedFieldsResponse) error {
	claims, err := wrappedClaims.UnwrapClaimsFromContext(r)
	if err != nil {
		log.Warn(err.Error())
		return err
	}

	updateAllowedFieldsResponse, err := a.administrator.UpdateAllowedFields(&administrator.UpdateAllowedFieldsRequest{
		Claims:  claims,
		Backend: request.Backend,
	})
	if err != nil {
		return err
	}

	response.Backend = updateAllowedFieldsResponse.Backend

	return nil
}

package client

import (
	"github.com/iot-my-world/brain/internal/log"
	"github.com/iot-my-world/brain/pkg/action"
	jsonRpcServiceProvider "github.com/iot-my-world/brain/pkg/api/jsonRpc/service/provider"
	"github.com/iot-my-world/brain/pkg/party/client"
	"github.com/iot-my-world/brain/pkg/party/client/validator"
	wrappedClaims "github.com/iot-my-world/brain/pkg/security/claims/wrapped"
	"github.com/iot-my-world/brain/pkg/validate/reasonInvalid"
	"net/http"
)

type adaptor struct {
	clientValidator validator.Validator
}

func New(clientValidator validator.Validator) *adaptor {
	return &adaptor{
		clientValidator: clientValidator,
	}
}

func (a *adaptor) Name() jsonRpcServiceProvider.Name {
	return jsonRpcServiceProvider.Name(validator.ServiceProvider)
}

func (a *adaptor) MethodRequiresAuthorization(string) bool {
	return true
}

type ValidateRequest struct {
	Client client.Client `json:"client"`
	Action action.Action `json:"action"`
}

type ValidateResponse struct {
	ReasonsInvalid []reasonInvalid.ReasonInvalid `json:"reasonsInvalid"`
}

func (a *adaptor) Validate(r *http.Request, request *ValidateRequest, response *ValidateResponse) error {
	claims, err := wrappedClaims.UnwrapClaimsFromContext(r)
	if err != nil {
		log.Warn(err.Error())
		return err
	}

	validateUserResponse, err := a.clientValidator.Validate(&validator.ValidateRequest{
		Claims: claims,
		Client: request.Client,
		Action: request.Action,
	})
	if err != nil {
		return err
	}

	response.ReasonsInvalid = validateUserResponse.ReasonsInvalid

	return nil
}

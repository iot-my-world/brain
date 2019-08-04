package jsonRpc

import (
	"github.com/iot-my-world/brain/internal/log"
	"github.com/iot-my-world/brain/pkg/action"
	jsonRpcServiceProvider "github.com/iot-my-world/brain/pkg/api/jsonRpc/service/provider"
	wrappedClaims "github.com/iot-my-world/brain/pkg/security/claims/wrapped"
	sigfoxBackendDataCallbackMessage "github.com/iot-my-world/brain/pkg/sigfox/backend/callback/data/message"
	"github.com/iot-my-world/brain/pkg/sigfox/backend/callback/data/message/validator"
	"github.com/iot-my-world/brain/pkg/validate/reasonInvalid"
	"net/http"
)

type adaptor struct {
	sigfoxBackendDataCallbackMessageValidator validator.Validator
}

func New(sigfoxBackendDataCallbackMessageValidator validator.Validator) *adaptor {
	return &adaptor{
		sigfoxBackendDataCallbackMessageValidator: sigfoxBackendDataCallbackMessageValidator,
	}
}

func (a *adaptor) Name() jsonRpcServiceProvider.Name {
	return jsonRpcServiceProvider.Name(validator.ServiceProvider)
}

func (a *adaptor) MethodRequiresAuthorization(string) bool {
	return true
}

type ValidateRequest struct {
	Message sigfoxBackendDataCallbackMessage.Message `json:"message"`
	Action  action.Action                            `json:"action"`
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

	validateMessageDeviceResponse, err := a.sigfoxBackendDataCallbackMessageValidator.Validate(&validator.ValidateRequest{
		Claims:  claims,
		Message: request.Message,
		Action:  request.Action,
	})
	if err != nil {
		return err
	}

	response.ReasonsInvalid = validateMessageDeviceResponse.ReasonsInvalid

	return nil
}

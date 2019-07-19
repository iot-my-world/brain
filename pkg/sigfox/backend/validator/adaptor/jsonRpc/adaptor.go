package jsonRpc

import (
	"github.com/iot-my-world/brain/internal/log"
	"github.com/iot-my-world/brain/pkg/action"
	jsonRpcServiceProvider "github.com/iot-my-world/brain/pkg/api/jsonRpc/service/provider"
	wrappedClaims "github.com/iot-my-world/brain/pkg/security/claims/wrapped"
	"github.com/iot-my-world/brain/pkg/sigfox/backend"
	"github.com/iot-my-world/brain/pkg/sigfox/backend/validator"
	"github.com/iot-my-world/brain/pkg/validate/reasonInvalid"
	"net/http"
)

type adaptor struct {
	backendDeviceValidator validator.Validator
}

func New(backendDeviceValidator validator.Validator) *adaptor {
	return &adaptor{
		backendDeviceValidator: backendDeviceValidator,
	}
}

func (a *adaptor) Name() jsonRpcServiceProvider.Name {
	return jsonRpcServiceProvider.Name(validator.ServiceProvider)
}

func (a *adaptor) MethodRequiresAuthorization(string) bool {
	return true
}

type ValidateRequest struct {
	Backend backend.Backend `json:"backend"`
	Action  action.Action   `json:"action"`
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

	validateBackendDeviceResponse, err := a.backendDeviceValidator.Validate(&validator.ValidateRequest{
		Claims:  claims,
		Backend: request.Backend,
		Action:  request.Action,
	})
	if err != nil {
		return err
	}

	response.ReasonsInvalid = validateBackendDeviceResponse.ReasonsInvalid

	return nil
}

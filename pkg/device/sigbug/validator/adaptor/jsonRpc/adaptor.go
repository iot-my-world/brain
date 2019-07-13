package jsonRpc

import (
	"github.com/iot-my-world/brain/internal/log"
	"github.com/iot-my-world/brain/pkg/action"
	jsonRpcServiceProvider "github.com/iot-my-world/brain/pkg/api/jsonRpc/service/provider"
	"github.com/iot-my-world/brain/pkg/device/sigbug"
	"github.com/iot-my-world/brain/pkg/device/sigbug/validator"
	wrappedClaims "github.com/iot-my-world/brain/pkg/security/claims/wrapped"
	"github.com/iot-my-world/brain/pkg/validate/reasonInvalid"
	"net/http"
)

type adaptor struct {
	sigbugDeviceValidator validator.Validator
}

func New(sigbugDeviceValidator validator.Validator) *adaptor {
	return &adaptor{
		sigbugDeviceValidator: sigbugDeviceValidator,
	}
}

func (a *adaptor) Name() jsonRpcServiceProvider.Name {
	return jsonRpcServiceProvider.Name(validator.ServiceProvider)
}

func (a *adaptor) MethodRequiresAuthorization(string) bool {
	return true
}

type ValidateRequest struct {
	Sigbug sigbug.Sigbug `json:"sigbug"`
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

	validateSigbugDeviceResponse, err := a.sigbugDeviceValidator.Validate(&validator.ValidateRequest{
		Claims: claims,
		Sigbug: request.Sigbug,
		Action: request.Action,
	})
	if err != nil {
		return err
	}

	response.ReasonsInvalid = validateSigbugDeviceResponse.ReasonsInvalid

	return nil
}

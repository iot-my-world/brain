package jsonRpc

import (
	"github.com/iot-my-world/brain/internal/log"
	"github.com/iot-my-world/brain/pkg/action"
	jsonRpcServiceProvider "github.com/iot-my-world/brain/pkg/api/jsonRpc/service/provider"
	sigbugGPSReading "github.com/iot-my-world/brain/pkg/device/sigbug/reading/gps"
	wrappedClaims "github.com/iot-my-world/brain/pkg/security/claims/wrapped"
	"github.com/iot-my-world/brain/pkg/sigfox/backend/callback/data/message/validator"
	"github.com/iot-my-world/brain/pkg/validate/reasonInvalid"
	"net/http"
)

type adaptor struct {
	sigfoxBackendDataCallbackReadingValidator validator.Validator
}

func New(sigfoxBackendDataCallbackReadingValidator validator.Validator) *adaptor {
	return &adaptor{
		sigfoxBackendDataCallbackReadingValidator: sigfoxBackendDataCallbackReadingValidator,
	}
}

func (a *adaptor) Name() jsonRpcServiceProvider.Name {
	return jsonRpcServiceProvider.Name(validator.ServiceProvider)
}

func (a *adaptor) MethodRequiresAuthorization(string) bool {
	return true
}

type ValidateRequest struct {
	Reading sigbugGPSReading.Reading `json:"reading"`
	Action  action.Action            `json:"action"`
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

	validateReadingDeviceResponse, err := a.sigfoxBackendDataCallbackReadingValidator.Validate(&validator.ValidateRequest{
		Claims:  claims,
		Reading: request.Reading,
		Action:  request.Action,
	})
	if err != nil {
		return err
	}

	response.ReasonsInvalid = validateReadingDeviceResponse.ReasonsInvalid

	return nil
}

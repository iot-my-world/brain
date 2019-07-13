package apiUser

import (
	"github.com/iot-my-world/brain/internal/log"
	"github.com/iot-my-world/brain/pkg/action"
	jsonRpcServiceProvider "github.com/iot-my-world/brain/pkg/api/jsonRpc/service/provider"
	wrappedClaims "github.com/iot-my-world/brain/pkg/security/claims/wrapped"
	"github.com/iot-my-world/brain/pkg/user/api"
	"github.com/iot-my-world/brain/pkg/user/api/validator"
	"github.com/iot-my-world/brain/pkg/validate/reasonInvalid"
	"net/http"
)

type adaptor struct {
	apiUserDeviceValidator validator.Validator
}

func New(companyValidator validator.Validator) *adaptor {
	return &adaptor{
		apiUserDeviceValidator: companyValidator,
	}
}

func (a *adaptor) Name() jsonRpcServiceProvider.Name {
	return jsonRpcServiceProvider.Name(validator.ServiceProvider)
}

func (a *adaptor) MethodRequiresAuthorization(string) bool {
	return true
}

type ValidateRequest struct {
	User   api.User      `json:"apiUser"`
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

	validateUserDeviceResponse, err := a.apiUserDeviceValidator.Validate(&validator.ValidateRequest{
		Claims: claims,
		User:   request.User,
		Action: request.Action,
	})
	if err != nil {
		return err
	}

	response.ReasonsInvalid = validateUserDeviceResponse.ReasonsInvalid

	return nil
}

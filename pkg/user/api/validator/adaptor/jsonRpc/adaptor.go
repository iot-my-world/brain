package apiUser

import (
	"github.com/iot-my-world/brain/internal/log"
	"github.com/iot-my-world/brain/pkg/action"
	"github.com/iot-my-world/brain/pkg/user/api"
	"github.com/iot-my-world/brain/pkg/user/api/validator"
	"github.com/iot-my-world/brain/pkg/validate/reasonInvalid"
	wrappedClaims "github.com/iot-my-world/brain/security/claims/wrapped"
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

type ValidateRequest struct {
	User   api.User      `json:"apiUser"`
	Action action.Action `json:"action"`
}

type ValidateResponse struct {
	ReasonsInvalid []reasonInvalid.ReasonInvalid `json:"reasonsInvalid"`
}

func (s *adaptor) Validate(r *http.Request, request *ValidateRequest, response *ValidateResponse) error {
	claims, err := wrappedClaims.UnwrapClaimsFromContext(r)
	if err != nil {
		log.Warn(err.Error())
		return err
	}

	validateUserDeviceResponse, err := s.apiUserDeviceValidator.Validate(&validator.ValidateRequest{
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

package apiUser

import (
	"github.com/iot-my-world/brain/action"
	"github.com/iot-my-world/brain/log"
	wrappedClaims "github.com/iot-my-world/brain/security/claims/wrapped"
	apiUser "github.com/iot-my-world/brain/user/api"
	apiUserDeviceValidator "github.com/iot-my-world/brain/user/api/validator"
	"github.com/iot-my-world/brain/validate/reasonInvalid"
	"net/http"
)

type adaptor struct {
	apiUserDeviceValidator apiUserDeviceValidator.Validator
}

func New(companyValidator apiUserDeviceValidator.Validator) *adaptor {
	return &adaptor{
		apiUserDeviceValidator: companyValidator,
	}
}

type ValidateRequest struct {
	User   apiUser.User  `json:"apiUser"`
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

	validateUserDeviceResponse, err := s.apiUserDeviceValidator.Validate(&apiUserDeviceValidator.ValidateRequest{
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

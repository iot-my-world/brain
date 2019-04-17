package apiUser

import (
	"gitlab.com/iotTracker/brain/action"
	"gitlab.com/iotTracker/brain/log"
	wrappedClaims "gitlab.com/iotTracker/brain/security/claims/wrapped"
	"gitlab.com/iotTracker/brain/tracker/device/apiUser"
	apiUserDeviceValidator "gitlab.com/iotTracker/brain/tracker/device/apiUser/validator"
	"gitlab.com/iotTracker/brain/validate/reasonInvalid"
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

package jsonRpc

import (
	"gitlab.com/iotTracker/brain/action"
	"gitlab.com/iotTracker/brain/log"
	wrappedClaims "gitlab.com/iotTracker/brain/security/claims/wrapped"
	"gitlab.com/iotTracker/brain/tracker/sf001"
	sf001DeviceValidator "gitlab.com/iotTracker/brain/tracker/sf001/validator"
	"gitlab.com/iotTracker/brain/validate/reasonInvalid"
	"net/http"
)

type adaptor struct {
	sf001DeviceValidator sf001DeviceValidator.Validator
}

func New(sf001DeviceValidator sf001DeviceValidator.Validator) *adaptor {
	return &adaptor{
		sf001DeviceValidator: sf001DeviceValidator,
	}
}

type ValidateRequest struct {
	SF001  sf001.SF001   `json:"sf001"`
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

	validateSF001DeviceResponse, err := s.sf001DeviceValidator.Validate(&sf001DeviceValidator.ValidateRequest{
		Claims: claims,
		SF001:  request.SF001,
		Action: request.Action,
	})
	if err != nil {
		return err
	}

	response.ReasonsInvalid = validateSF001DeviceResponse.ReasonsInvalid

	return nil
}

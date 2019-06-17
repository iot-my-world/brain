package jsonRpc

import (
	"github.com/iot-my-world/brain/action"
	"github.com/iot-my-world/brain/log"
	wrappedClaims "github.com/iot-my-world/brain/security/claims/wrapped"
	"github.com/iot-my-world/brain/tracker/sf001"
	sf001DeviceValidator "github.com/iot-my-world/brain/tracker/sf001/validator"
	"github.com/iot-my-world/brain/validate/reasonInvalid"
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

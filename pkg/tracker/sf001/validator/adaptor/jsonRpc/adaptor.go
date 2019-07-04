package jsonRpc

import (
	"github.com/iot-my-world/brain/log"
	"github.com/iot-my-world/brain/pkg/action"
	sf0012 "github.com/iot-my-world/brain/pkg/tracker/sf001"
	"github.com/iot-my-world/brain/pkg/tracker/sf001/validator"
	wrappedClaims "github.com/iot-my-world/brain/security/claims/wrapped"
	"github.com/iot-my-world/brain/validate/reasonInvalid"
	"net/http"
)

type adaptor struct {
	sf001DeviceValidator validator.Validator
}

func New(sf001DeviceValidator validator.Validator) *adaptor {
	return &adaptor{
		sf001DeviceValidator: sf001DeviceValidator,
	}
}

type ValidateRequest struct {
	SF001  sf0012.SF001  `json:"sf001"`
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

	validateSF001DeviceResponse, err := s.sf001DeviceValidator.Validate(&validator.ValidateRequest{
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

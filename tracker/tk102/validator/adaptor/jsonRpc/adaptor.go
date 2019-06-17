package tk102

import (
	"github.com/iot-my-world/brain/action"
	"github.com/iot-my-world/brain/log"
	wrappedClaims "github.com/iot-my-world/brain/security/claims/wrapped"
	"github.com/iot-my-world/brain/tracker/tk102"
	tk102DeviceValidator "github.com/iot-my-world/brain/tracker/tk102/validator"
	"github.com/iot-my-world/brain/validate/reasonInvalid"
	"net/http"
)

type adaptor struct {
	tk102DeviceValidator tk102DeviceValidator.Validator
}

func New(companyValidator tk102DeviceValidator.Validator) *adaptor {
	return &adaptor{
		tk102DeviceValidator: companyValidator,
	}
}

type ValidateRequest struct {
	TK102  tk102.TK102   `json:"tk102"`
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

	validateTK102DeviceResponse, err := s.tk102DeviceValidator.Validate(&tk102DeviceValidator.ValidateRequest{
		Claims: claims,
		TK102:  request.TK102,
		Action: request.Action,
	})
	if err != nil {
		return err
	}

	response.ReasonsInvalid = validateTK102DeviceResponse.ReasonsInvalid

	return nil
}

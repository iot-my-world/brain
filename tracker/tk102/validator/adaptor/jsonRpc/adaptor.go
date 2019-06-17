package tk102

import (
	"gitlab.com/iotTracker/brain/action"
	"gitlab.com/iotTracker/brain/log"
	wrappedClaims "gitlab.com/iotTracker/brain/security/claims/wrapped"
	"gitlab.com/iotTracker/brain/tracker/tk102"
	tk102DeviceValidator "gitlab.com/iotTracker/brain/tracker/tk102/validator"
	"gitlab.com/iotTracker/brain/validate/reasonInvalid"
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

package tk102

import (
	"github.com/iot-my-world/brain/internal/log"
	"github.com/iot-my-world/brain/pkg/action"
	tk1022 "github.com/iot-my-world/brain/pkg/tracker/tk102"
	"github.com/iot-my-world/brain/pkg/tracker/tk102/validator"
	"github.com/iot-my-world/brain/pkg/validate/reasonInvalid"
	wrappedClaims "github.com/iot-my-world/brain/security/claims/wrapped"
	"net/http"
)

type adaptor struct {
	tk102DeviceValidator validator.Validator
}

func New(companyValidator validator.Validator) *adaptor {
	return &adaptor{
		tk102DeviceValidator: companyValidator,
	}
}

type ValidateRequest struct {
	TK102  tk1022.TK102  `json:"tk102"`
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

	validateTK102DeviceResponse, err := s.tk102DeviceValidator.Validate(&validator.ValidateRequest{
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

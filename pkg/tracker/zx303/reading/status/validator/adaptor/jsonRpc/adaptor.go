package jsonRpc

import (
	"github.com/iot-my-world/brain/action"
	"github.com/iot-my-world/brain/log"
	"github.com/iot-my-world/brain/pkg/tracker/zx303/reading/status"
	"github.com/iot-my-world/brain/pkg/tracker/zx303/reading/status/validator"
	wrappedClaims "github.com/iot-my-world/brain/security/claims/wrapped"
	"github.com/iot-my-world/brain/validate/reasonInvalid"
	"net/http"
)

type adaptor struct {
	zx303StatusReadingValidator validator.Validator
}

func New(zx303StatusReadingValidator validator.Validator) *adaptor {
	return &adaptor{
		zx303StatusReadingValidator: zx303StatusReadingValidator,
	}
}

type ValidateRequest struct {
	ZX303StatusReading status.Reading `json:"zx303StatusReading"`
	Action             action.Action  `json:"action"`
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

	validateZX303StatusReadingResponse, err := s.zx303StatusReadingValidator.Validate(&validator.ValidateRequest{
		Claims:             claims,
		ZX303StatusReading: request.ZX303StatusReading,
		Action:             request.Action,
	})
	if err != nil {
		return err
	}

	response.ReasonsInvalid = validateZX303StatusReadingResponse.ReasonsInvalid

	return nil
}

package jsonRpc

import (
	"github.com/iot-my-world/brain/action"
	"github.com/iot-my-world/brain/log"
	"github.com/iot-my-world/brain/pkg/tracker/zx303/reading/gps"
	"github.com/iot-my-world/brain/pkg/tracker/zx303/reading/gps/validator"
	wrappedClaims "github.com/iot-my-world/brain/security/claims/wrapped"
	"github.com/iot-my-world/brain/validate/reasonInvalid"
	"net/http"
)

type adaptor struct {
	zx303GPSReadingValidator validator.Validator
}

func New(zx303GPSReadingValidator validator.Validator) *adaptor {
	return &adaptor{
		zx303GPSReadingValidator: zx303GPSReadingValidator,
	}
}

type ValidateRequest struct {
	ZX303GPSReading gps.Reading   `json:"zx303GPSReading"`
	Action          action.Action `json:"action"`
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

	validateZX303GPSReadingResponse, err := s.zx303GPSReadingValidator.Validate(&validator.ValidateRequest{
		Claims:          claims,
		ZX303GPSReading: request.ZX303GPSReading,
		Action:          request.Action,
	})
	if err != nil {
		return err
	}

	response.ReasonsInvalid = validateZX303GPSReadingResponse.ReasonsInvalid

	return nil
}

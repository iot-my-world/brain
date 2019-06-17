package jsonRpc

import (
	"gitlab.com/iotTracker/brain/action"
	"gitlab.com/iotTracker/brain/log"
	wrappedClaims "gitlab.com/iotTracker/brain/security/claims/wrapped"
	zx303GPSReading "gitlab.com/iotTracker/brain/tracker/zx303/reading/gps"
	zx303GPSReadingValidator "gitlab.com/iotTracker/brain/tracker/zx303/reading/gps/validator"
	"gitlab.com/iotTracker/brain/validate/reasonInvalid"
	"net/http"
)

type adaptor struct {
	zx303GPSReadingValidator zx303GPSReadingValidator.Validator
}

func New(zx303GPSReadingValidator zx303GPSReadingValidator.Validator) *adaptor {
	return &adaptor{
		zx303GPSReadingValidator: zx303GPSReadingValidator,
	}
}

type ValidateRequest struct {
	ZX303GPSReading zx303GPSReading.Reading `json:"zx303GPSReading"`
	Action          action.Action           `json:"action"`
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

	validateZX303GPSReadingResponse, err := s.zx303GPSReadingValidator.Validate(&zx303GPSReadingValidator.ValidateRequest{
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

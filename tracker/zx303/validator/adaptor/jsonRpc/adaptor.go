package zx303

import (
	"gitlab.com/iotTracker/brain/action"
	"gitlab.com/iotTracker/brain/log"
	wrappedClaims "gitlab.com/iotTracker/brain/security/claims/wrapped"
	"gitlab.com/iotTracker/brain/tracker/zx303"
	zx303DeviceValidator "gitlab.com/iotTracker/brain/tracker/zx303/validator"
	"gitlab.com/iotTracker/brain/validate/reasonInvalid"
	"net/http"
)

type adaptor struct {
	zx303DeviceValidator zx303DeviceValidator.Validator
}

func New(companyValidator zx303DeviceValidator.Validator) *adaptor {
	return &adaptor{
		zx303DeviceValidator: companyValidator,
	}
}

type ValidateRequest struct {
	ZX303  zx303.ZX303   `json:"zx303"`
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

	validateZX303DeviceResponse, err := s.zx303DeviceValidator.Validate(&zx303DeviceValidator.ValidateRequest{
		Claims: claims,
		ZX303:  request.ZX303,
		Action: request.Action,
	})
	if err != nil {
		return err
	}

	response.ReasonsInvalid = validateZX303DeviceResponse.ReasonsInvalid

	return nil
}

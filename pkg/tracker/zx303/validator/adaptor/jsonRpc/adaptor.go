package jsonRpc

import (
	"github.com/iot-my-world/brain/internal/log"
	"github.com/iot-my-world/brain/pkg/action"
	wrappedClaims "github.com/iot-my-world/brain/pkg/security/claims/wrapped"
	zx3032 "github.com/iot-my-world/brain/pkg/tracker/zx303"
	"github.com/iot-my-world/brain/pkg/tracker/zx303/validator"
	"github.com/iot-my-world/brain/pkg/validate/reasonInvalid"
	"net/http"
)

type adaptor struct {
	zx303DeviceValidator validator.Validator
}

func New(zx303DeviceValidator validator.Validator) *adaptor {
	return &adaptor{
		zx303DeviceValidator: zx303DeviceValidator,
	}
}

type ValidateRequest struct {
	ZX303  zx3032.ZX303  `json:"zx303"`
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

	validateZX303DeviceResponse, err := s.zx303DeviceValidator.Validate(&validator.ValidateRequest{
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

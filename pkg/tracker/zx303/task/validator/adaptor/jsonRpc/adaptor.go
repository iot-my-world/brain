package jsonRpc

import (
	"github.com/iot-my-world/brain/action"
	"github.com/iot-my-world/brain/log"
	"github.com/iot-my-world/brain/pkg/tracker/zx303/task"
	"github.com/iot-my-world/brain/pkg/tracker/zx303/task/validator"
	wrappedClaims "github.com/iot-my-world/brain/security/claims/wrapped"
	"github.com/iot-my-world/brain/validate/reasonInvalid"
	"net/http"
)

type adaptor struct {
	zx303TaskValidator validator.Validator
}

func New(zx303TaskValidator validator.Validator) *adaptor {
	return &adaptor{
		zx303TaskValidator: zx303TaskValidator,
	}
}

type ValidateRequest struct {
	ZX303Task task.Task     `json:"zx303Task"`
	Action    action.Action `json:"action"`
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

	validateZX303TaskResponse, err := s.zx303TaskValidator.Validate(&validator.ValidateRequest{
		Claims:    claims,
		ZX303Task: request.ZX303Task,
		Action:    request.Action,
	})
	if err != nil {
		return err
	}

	response.ReasonsInvalid = validateZX303TaskResponse.ReasonsInvalid

	return nil
}

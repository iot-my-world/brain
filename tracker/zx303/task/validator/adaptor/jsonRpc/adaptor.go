package jsonRpc

import (
	"gitlab.com/iotTracker/brain/action"
	"gitlab.com/iotTracker/brain/log"
	wrappedClaims "gitlab.com/iotTracker/brain/security/claims/wrapped"
	zx303Task "gitlab.com/iotTracker/brain/tracker/zx303/task"
	zx303TaskValidator "gitlab.com/iotTracker/brain/tracker/zx303/task/validator"
	"gitlab.com/iotTracker/brain/validate/reasonInvalid"
	"net/http"
)

type adaptor struct {
	zx303TaskValidator zx303TaskValidator.Validator
}

func New(zx303TaskValidator zx303TaskValidator.Validator) *adaptor {
	return &adaptor{
		zx303TaskValidator: zx303TaskValidator,
	}
}

type ValidateRequest struct {
	ZX303Task zx303Task.Task `json:"zx303Task"`
	Action    action.Action  `json:"action"`
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

	validateZX303TaskResponse, err := s.zx303TaskValidator.Validate(&zx303TaskValidator.ValidateRequest{
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
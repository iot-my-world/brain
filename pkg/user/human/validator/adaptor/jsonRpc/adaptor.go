package user

import (
	"github.com/iot-my-world/brain/action"
	"github.com/iot-my-world/brain/log"
	"github.com/iot-my-world/brain/pkg/user/human"
	"github.com/iot-my-world/brain/pkg/user/human/validator"
	wrappedClaims "github.com/iot-my-world/brain/security/claims/wrapped"
	"github.com/iot-my-world/brain/validate/reasonInvalid"
	"net/http"
)

type adaptor struct {
	RecordHandler validator.Validator
}

func New(recordHandler validator.Validator) *adaptor {
	return &adaptor{
		RecordHandler: recordHandler,
	}
}

type ValidateRequest struct {
	User   human.User    `json:"user"`
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

	validateUserResponse, err := s.RecordHandler.Validate(&validator.ValidateRequest{
		Claims: claims,
		User:   request.User,
		Action: request.Action,
	})
	if err != nil {
		return err
	}

	response.ReasonsInvalid = validateUserResponse.ReasonsInvalid

	return nil
}

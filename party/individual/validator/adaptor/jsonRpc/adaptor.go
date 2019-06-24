package jsonRpc

import (
	"github.com/iot-my-world/brain/action"
	"github.com/iot-my-world/brain/log"
	"github.com/iot-my-world/brain/party/individual"
	individualIndividualValidator "github.com/iot-my-world/brain/party/individual/validator"
	wrappedClaims "github.com/iot-my-world/brain/security/claims/wrapped"
	"github.com/iot-my-world/brain/validate/reasonInvalid"
	"net/http"
)

type adaptor struct {
	individualIndividualValidator individualIndividualValidator.Validator
}

func New(individualIndividualValidator individualIndividualValidator.Validator) *adaptor {
	return &adaptor{
		individualIndividualValidator: individualIndividualValidator,
	}
}

type ValidateRequest struct {
	Individual individual.Individual `json:"individual"`
	Action     action.Action         `json:"action"`
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

	validateIndividualIndividualResponse, err := s.individualIndividualValidator.Validate(&individualIndividualValidator.ValidateRequest{
		Claims:     claims,
		Individual: request.Individual,
		Action:     request.Action,
	})
	if err != nil {
		return err
	}

	response.ReasonsInvalid = validateIndividualIndividualResponse.ReasonsInvalid

	return nil
}

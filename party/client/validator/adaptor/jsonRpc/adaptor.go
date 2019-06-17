package client

import (
	"gitlab.com/iotTracker/brain/action"
	"gitlab.com/iotTracker/brain/log"
	"gitlab.com/iotTracker/brain/party/client"
	clientValidator "gitlab.com/iotTracker/brain/party/client/validator"
	wrappedClaims "gitlab.com/iotTracker/brain/security/claims/wrapped"
	"gitlab.com/iotTracker/brain/validate/reasonInvalid"
	"net/http"
)

type adaptor struct {
	clientValidator clientValidator.Validator
}

func New(clientValidator clientValidator.Validator) *adaptor {
	return &adaptor{
		clientValidator: clientValidator,
	}
}

type ValidateRequest struct {
	Client client.Client `json:"client"`
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

	validateUserResponse, err := s.clientValidator.Validate(&clientValidator.ValidateRequest{
		Claims: claims,
		Client: request.Client,
		Action: request.Action,
	})
	if err != nil {
		return err
	}

	response.ReasonsInvalid = validateUserResponse.ReasonsInvalid

	return nil
}

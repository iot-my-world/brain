package client

import (
	"github.com/iot-my-world/brain/internal/log"
	"github.com/iot-my-world/brain/pkg/action"
	client2 "github.com/iot-my-world/brain/pkg/party/client"
	"github.com/iot-my-world/brain/pkg/party/client/validator"
	"github.com/iot-my-world/brain/pkg/validate/reasonInvalid"
	wrappedClaims "github.com/iot-my-world/brain/security/claims/wrapped"
	"net/http"
)

type adaptor struct {
	clientValidator validator.Validator
}

func New(clientValidator validator.Validator) *adaptor {
	return &adaptor{
		clientValidator: clientValidator,
	}
}

type ValidateRequest struct {
	Client client2.Client `json:"client"`
	Action action.Action  `json:"action"`
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

	validateUserResponse, err := s.clientValidator.Validate(&validator.ValidateRequest{
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

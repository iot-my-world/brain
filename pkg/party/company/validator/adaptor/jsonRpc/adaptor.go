package adaptor

import (
	"github.com/iot-my-world/brain/internal/log"
	"github.com/iot-my-world/brain/pkg/action"
	"github.com/iot-my-world/brain/pkg/party/company"
	companyValidator "github.com/iot-my-world/brain/pkg/party/company/validator"
	wrappedClaims "github.com/iot-my-world/brain/pkg/security/claims/wrapped"
	"github.com/iot-my-world/brain/pkg/validate/reasonInvalid"
	"net/http"
)

type adaptor struct {
	companyValidator companyValidator.Validator
}

func New(companyValidator companyValidator.Validator) *adaptor {
	return &adaptor{
		companyValidator: companyValidator,
	}
}

type ValidateRequest struct {
	Company company.Company `json:"company"`
	Action  action.Action   `json:"action"`
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

	validateUserResponse, err := s.companyValidator.Validate(&companyValidator.ValidateRequest{
		Claims:  claims,
		Company: request.Company,
		Action:  request.Action,
	})
	if err != nil {
		return err
	}

	response.ReasonsInvalid = validateUserResponse.ReasonsInvalid

	return nil
}

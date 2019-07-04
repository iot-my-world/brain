package company

import (
	"github.com/iot-my-world/brain/log"
	"github.com/iot-my-world/brain/pkg/action"
	company2 "github.com/iot-my-world/brain/pkg/party/company"
	"github.com/iot-my-world/brain/pkg/party/company/validator"
	wrappedClaims "github.com/iot-my-world/brain/security/claims/wrapped"
	"github.com/iot-my-world/brain/validate/reasonInvalid"
	"net/http"
)

type adaptor struct {
	companyValidator validator.Validator
}

func New(companyValidator validator.Validator) *adaptor {
	return &adaptor{
		companyValidator: companyValidator,
	}
}

type ValidateRequest struct {
	Company company2.Company `json:"company"`
	Action  action.Action    `json:"action"`
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

	validateUserResponse, err := s.companyValidator.Validate(&validator.ValidateRequest{
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

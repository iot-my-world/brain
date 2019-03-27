package company

import (
	"gitlab.com/iotTracker/brain/action"
	"gitlab.com/iotTracker/brain/log"
	"gitlab.com/iotTracker/brain/party/company"
	companyValidator "gitlab.com/iotTracker/brain/party/company/validator"
	wrappedClaims "gitlab.com/iotTracker/brain/security/claims/wrapped"
	"gitlab.com/iotTracker/brain/validate/reasonInvalid"
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

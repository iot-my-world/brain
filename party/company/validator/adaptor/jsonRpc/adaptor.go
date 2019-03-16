package company

import (
	"gitlab.com/iotTracker/brain/action"
	"gitlab.com/iotTracker/brain/log"
	"gitlab.com/iotTracker/brain/party/company"
	userValidator "gitlab.com/iotTracker/brain/party/company/validator"
	"gitlab.com/iotTracker/brain/security/wrappedClaims"
	"gitlab.com/iotTracker/brain/validate/reasonInvalid"
	"net/http"
)

type adaptor struct {
	RecordHandler userValidator.Validator
}

func New(recordHandler userValidator.Validator) *adaptor {
	return &adaptor{
		RecordHandler: recordHandler,
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

	validateUserResponse := userValidator.ValidateResponse{}
	if err := s.RecordHandler.Validate(&userValidator.ValidateRequest{
		Claims:  claims,
		Company: request.Company,
		Action:  request.Action,
	}, &validateUserResponse); err != nil {
		return err
	}

	response.ReasonsInvalid = validateUserResponse.ReasonsInvalid

	return nil
}

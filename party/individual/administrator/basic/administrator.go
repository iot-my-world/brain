package basic

import (
	"fmt"
	brainException "github.com/iot-my-world/brain/exception"
	individualIndividualAction "github.com/iot-my-world/brain/party/individual/action"
	individualIndividualAdministrator "github.com/iot-my-world/brain/party/individual/administrator"
	individualIndividualAdministratorException "github.com/iot-my-world/brain/party/individual/administrator/exception"
	individualRecordHandler "github.com/iot-my-world/brain/party/individual/recordHandler"
	individualIndividualValidator "github.com/iot-my-world/brain/party/individual/validator"
	"github.com/iot-my-world/brain/search/identifier/id"
)

type administrator struct {
	individualIndividualValidator individualIndividualValidator.Validator
	individualRecordHandler       *individualRecordHandler.RecordHandler
}

func New(
	individualIndividualValidator individualIndividualValidator.Validator,
	individualRecordHandler *individualRecordHandler.RecordHandler,
) individualIndividualAdministrator.Administrator {
	return &administrator{
		individualIndividualValidator: individualIndividualValidator,
		individualRecordHandler:       individualRecordHandler,
	}
}

func (a *administrator) ValidateCreateRequest(request *individualIndividualAdministrator.CreateRequest) error {
	reasonsInvalid := make([]string, 0)

	if request.Claims == nil {
		reasonsInvalid = append(reasonsInvalid, "claims are nil")
	} else {
		individualIndividualValidateResponse, err := a.individualIndividualValidator.Validate(&individualIndividualValidator.ValidateRequest{
			Claims:     request.Claims,
			Individual: request.Individual,
			Action:     individualIndividualAction.Create,
		})
		if err != nil {
			reasonsInvalid = append(reasonsInvalid, "error validating individual individual: "+err.Error())
		}
		if len(individualIndividualValidateResponse.ReasonsInvalid) > 0 {
			for _, reason := range individualIndividualValidateResponse.ReasonsInvalid {
				reasonsInvalid = append(reasonsInvalid, fmt.Sprintf("individual individual invalid: %s - %s - %s", reason.Field, reason.Type, reason.Help))
			}
		}
	}

	if len(reasonsInvalid) > 0 {
		return brainException.RequestInvalid{Reasons: reasonsInvalid}
	}

	return nil
}

func (a *administrator) Create(request *individualIndividualAdministrator.CreateRequest) (*individualIndividualAdministrator.CreateResponse, error) {
	if err := a.ValidateCreateRequest(request); err != nil {
		return nil, err
	}

	createResponse, err := a.individualRecordHandler.Create(&individualRecordHandler.CreateRequest{
		Individual: request.Individual,
	})
	if err != nil {
		return nil, individualIndividualAdministratorException.IndividualCreation{Reasons: []string{err.Error()}}
	}

	return &individualIndividualAdministrator.CreateResponse{
		Individual: createResponse.Individual,
	}, nil
}

func (a *administrator) ValidateUpdateAllowedFieldsRequest(request *individualIndividualAdministrator.UpdateAllowedFieldsRequest) error {
	reasonsInvalid := make([]string, 0)

	if request.Claims == nil {
		reasonsInvalid = append(reasonsInvalid, "claims are nil")
	} else {
		// individual must be valid
		validationResponse, err := a.individualIndividualValidator.Validate(&individualIndividualValidator.ValidateRequest{
			Claims: request.Claims,
			Action: individualIndividualAction.UpdateAllowedFields,
		})
		if err != nil {
			reasonsInvalid = append(reasonsInvalid, "error validating individual: "+err.Error())
		}
		if len(validationResponse.ReasonsInvalid) > 0 {
			for _, reason := range validationResponse.ReasonsInvalid {
				reasonsInvalid = append(reasonsInvalid, fmt.Sprintf("individual individual invalid: %s - %s - %s", reason.Field, reason.Type, reason.Help))
			}
		}
	}

	if len(reasonsInvalid) > 0 {
		return brainException.RequestInvalid{Reasons: reasonsInvalid}
	}
	return nil
}

func (a *administrator) UpdateAllowedFields(request *individualIndividualAdministrator.UpdateAllowedFieldsRequest) (*individualIndividualAdministrator.UpdateAllowedFieldsResponse, error) {
	if err := a.ValidateUpdateAllowedFieldsRequest(request); err != nil {
		return nil, err
	}

	// retrieve the individual
	individualRetrieveResponse, err := a.individualRecordHandler.Retrieve(&individualRecordHandler.RetrieveRequest{
		Claims:     request.Claims,
		Identifier: id.Identifier{Id: request.Individual.Id},
	})
	if err != nil {
		return nil, individualIndividualAdministratorException.IndividualRetrieval{Reasons: []string{err.Error()}}
	}

	// update the allowed fields on the individual

	// update the individual
	_, err = a.individualRecordHandler.Update(&individualRecordHandler.UpdateRequest{
		Claims:     request.Claims,
		Identifier: id.Identifier{Id: request.Individual.Id},
		Individual: individualRetrieveResponse.Individual,
	})
	if err != nil {
		return nil, individualIndividualAdministratorException.IndividualUpdate{Reasons: []string{err.Error()}}
	}

	return &individualIndividualAdministrator.UpdateAllowedFieldsResponse{
		Individual: individualRetrieveResponse.Individual,
	}, nil
}

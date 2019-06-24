package basic

import (
	"fmt"
	brainException "github.com/iot-my-world/brain/exception"
	sf001DeviceAction "github.com/iot-my-world/brain/party/individual/action"
	sf001DeviceAdministrator "github.com/iot-my-world/brain/party/individual/administrator"
	sf001DeviceAdministratorException "github.com/iot-my-world/brain/party/individual/administrator/exception"
	sf001RecordHandler "github.com/iot-my-world/brain/party/individual/recordHandler"
	sf001DeviceValidator "github.com/iot-my-world/brain/party/individual/validator"
	"github.com/iot-my-world/brain/search/identifier/id"
)

type administrator struct {
	sf001DeviceValidator sf001DeviceValidator.Validator
	sf001RecordHandler   *sf001RecordHandler.RecordHandler
}

func New(
	sf001DeviceValidator sf001DeviceValidator.Validator,
	sf001RecordHandler *sf001RecordHandler.RecordHandler,
) sf001DeviceAdministrator.Administrator {
	return &administrator{
		sf001DeviceValidator: sf001DeviceValidator,
		sf001RecordHandler:   sf001RecordHandler,
	}
}

func (a *administrator) ValidateCreateRequest(request *sf001DeviceAdministrator.CreateRequest) error {
	reasonsInvalid := make([]string, 0)

	if request.Claims == nil {
		reasonsInvalid = append(reasonsInvalid, "claims are nil")
	} else {
		sf001DeviceValidateResponse, err := a.sf001DeviceValidator.Validate(&sf001DeviceValidator.ValidateRequest{
			Claims:     request.Claims,
			Individual: request.Individual,
			Action:     sf001DeviceAction.Create,
		})
		if err != nil {
			reasonsInvalid = append(reasonsInvalid, "error validating individual device: "+err.Error())
		}
		if len(sf001DeviceValidateResponse.ReasonsInvalid) > 0 {
			for _, reason := range sf001DeviceValidateResponse.ReasonsInvalid {
				reasonsInvalid = append(reasonsInvalid, fmt.Sprintf("individual device invalid: %s - %s - %s", reason.Field, reason.Type, reason.Help))
			}
		}
	}

	if len(reasonsInvalid) > 0 {
		return brainException.RequestInvalid{Reasons: reasonsInvalid}
	}

	return nil
}

func (a *administrator) Create(request *sf001DeviceAdministrator.CreateRequest) (*sf001DeviceAdministrator.CreateResponse, error) {
	if err := a.ValidateCreateRequest(request); err != nil {
		return nil, err
	}

	createResponse, err := a.sf001RecordHandler.Create(&sf001RecordHandler.CreateRequest{
		Individual: request.Individual,
	})
	if err != nil {
		return nil, sf001DeviceAdministratorException.DeviceCreation{Reasons: []string{err.Error()}}
	}

	return &sf001DeviceAdministrator.CreateResponse{
		Individual: createResponse.Individual,
	}, nil
}

func (a *administrator) ValidateUpdateAllowedFieldsRequest(request *sf001DeviceAdministrator.UpdateAllowedFieldsRequest) error {
	reasonsInvalid := make([]string, 0)

	if request.Claims == nil {
		reasonsInvalid = append(reasonsInvalid, "claims are nil")
	} else {
		// device must be valid
		validationResponse, err := a.sf001DeviceValidator.Validate(&sf001DeviceValidator.ValidateRequest{
			Claims: request.Claims,
			Action: sf001DeviceAction.UpdateAllowedFields,
		})
		if err != nil {
			reasonsInvalid = append(reasonsInvalid, "error validating device: "+err.Error())
		}
		if len(validationResponse.ReasonsInvalid) > 0 {
			for _, reason := range validationResponse.ReasonsInvalid {
				reasonsInvalid = append(reasonsInvalid, fmt.Sprintf("individual device invalid: %s - %s - %s", reason.Field, reason.Type, reason.Help))
			}
		}
	}

	if len(reasonsInvalid) > 0 {
		return brainException.RequestInvalid{Reasons: reasonsInvalid}
	}
	return nil
}

func (a *administrator) UpdateAllowedFields(request *sf001DeviceAdministrator.UpdateAllowedFieldsRequest) (*sf001DeviceAdministrator.UpdateAllowedFieldsResponse, error) {
	if err := a.ValidateUpdateAllowedFieldsRequest(request); err != nil {
		return nil, err
	}

	// retrieve the device
	deviceRetrieveResponse, err := a.sf001RecordHandler.Retrieve(&sf001RecordHandler.RetrieveRequest{
		Claims:     request.Claims,
		Identifier: id.Identifier{Id: request.Individual.Id},
	})
	if err != nil {
		return nil, sf001DeviceAdministratorException.DeviceRetrieval{Reasons: []string{err.Error()}}
	}

	// update the allowed fields on the device

	// update the device
	_, err = a.sf001RecordHandler.Update(&sf001RecordHandler.UpdateRequest{
		Claims:     request.Claims,
		Identifier: id.Identifier{Id: request.Individual.Id},
		Individual: deviceRetrieveResponse.Individual,
	})
	if err != nil {
		return nil, sf001DeviceAdministratorException.DeviceUpdate{Reasons: []string{err.Error()}}
	}

	return &sf001DeviceAdministrator.UpdateAllowedFieldsResponse{
		Individual: deviceRetrieveResponse.Individual,
	}, nil
}

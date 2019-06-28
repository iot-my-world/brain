package basic

import (
	"fmt"
	brainException "github.com/iot-my-world/brain/exception"
	"github.com/iot-my-world/brain/search/identifier/id"
	sf001DeviceAction "github.com/iot-my-world/brain/tracker/sf001/action"
	sf001DeviceAdministrator "github.com/iot-my-world/brain/tracker/sf001/administrator"
	sf001DeviceAdministratorException "github.com/iot-my-world/brain/tracker/sf001/administrator/exception"
	sf001RecordHandler "github.com/iot-my-world/brain/tracker/sf001/recordHandler"
	sf001DeviceValidator "github.com/iot-my-world/brain/tracker/sf001/validator"
)

type administrator struct {
	sf001DeviceValidator sf001DeviceValidator.Validator
	sf001RecordHandler   sf001RecordHandler.RecordHandler
}

func New(
	sf001DeviceValidator sf001DeviceValidator.Validator,
	sf001RecordHandler sf001RecordHandler.RecordHandler,
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
			Claims: request.Claims,
			SF001:  request.SF001,
			Action: sf001DeviceAction.Create,
		})
		if err != nil {
			reasonsInvalid = append(reasonsInvalid, "error validating sf001 device: "+err.Error())
		} else {
			if len(sf001DeviceValidateResponse.ReasonsInvalid) > 0 {
				for _, reason := range sf001DeviceValidateResponse.ReasonsInvalid {
					reasonsInvalid = append(reasonsInvalid, fmt.Sprintf("sf001 device invalid: %s - %s - %s", reason.Field, reason.Type, reason.Help))
				}
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
		SF001: request.SF001,
	})
	if err != nil {
		return nil, sf001DeviceAdministratorException.DeviceCreation{Reasons: []string{err.Error()}}
	}

	return &sf001DeviceAdministrator.CreateResponse{
		SF001: createResponse.SF001,
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
		} else {
			if len(validationResponse.ReasonsInvalid) > 0 {
				for _, reason := range validationResponse.ReasonsInvalid {
					reasonsInvalid = append(reasonsInvalid, fmt.Sprintf("sf001 device invalid: %s - %s - %s", reason.Field, reason.Type, reason.Help))
				}
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
		Identifier: id.Identifier{Id: request.SF001.Id},
	})
	if err != nil {
		return nil, sf001DeviceAdministratorException.DeviceRetrieval{Reasons: []string{err.Error()}}
	}

	// update the allowed fields on the device

	// update the device
	_, err = a.sf001RecordHandler.Update(&sf001RecordHandler.UpdateRequest{
		Claims:     request.Claims,
		Identifier: id.Identifier{Id: request.SF001.Id},
		SF001:      deviceRetrieveResponse.SF001,
	})
	if err != nil {
		return nil, sf001DeviceAdministratorException.DeviceUpdate{Reasons: []string{err.Error()}}
	}

	return &sf001DeviceAdministrator.UpdateAllowedFieldsResponse{
		SF001: deviceRetrieveResponse.SF001,
	}, nil
}

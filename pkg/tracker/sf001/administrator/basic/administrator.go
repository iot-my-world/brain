package basic

import (
	"fmt"
	brainException "github.com/iot-my-world/brain/internal/exception"
	"github.com/iot-my-world/brain/pkg/search/identifier/id"
	"github.com/iot-my-world/brain/pkg/tracker/sf001/action"
	administrator2 "github.com/iot-my-world/brain/pkg/tracker/sf001/administrator"
	"github.com/iot-my-world/brain/pkg/tracker/sf001/administrator/exception"
	"github.com/iot-my-world/brain/pkg/tracker/sf001/recordHandler"
	"github.com/iot-my-world/brain/pkg/tracker/sf001/validator"
)

type administrator struct {
	sf001DeviceValidator validator.Validator
	sf001RecordHandler   recordHandler.RecordHandler
}

func New(
	sf001DeviceValidator validator.Validator,
	sf001RecordHandler recordHandler.RecordHandler,
) administrator2.Administrator {
	return &administrator{
		sf001DeviceValidator: sf001DeviceValidator,
		sf001RecordHandler:   sf001RecordHandler,
	}
}

func (a *administrator) ValidateCreateRequest(request *administrator2.CreateRequest) error {
	reasonsInvalid := make([]string, 0)

	if request.Claims == nil {
		reasonsInvalid = append(reasonsInvalid, "claims are nil")
	} else {
		sf001DeviceValidateResponse, err := a.sf001DeviceValidator.Validate(&validator.ValidateRequest{
			Claims: request.Claims,
			SF001:  request.SF001,
			Action: action.Create,
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

func (a *administrator) Create(request *administrator2.CreateRequest) (*administrator2.CreateResponse, error) {
	if err := a.ValidateCreateRequest(request); err != nil {
		return nil, err
	}

	createResponse, err := a.sf001RecordHandler.Create(&recordHandler.CreateRequest{
		SF001: request.SF001,
	})
	if err != nil {
		return nil, exception.DeviceCreation{Reasons: []string{err.Error()}}
	}

	return &administrator2.CreateResponse{
		SF001: createResponse.SF001,
	}, nil
}

func (a *administrator) ValidateUpdateAllowedFieldsRequest(request *administrator2.UpdateAllowedFieldsRequest) error {
	reasonsInvalid := make([]string, 0)

	if request.Claims == nil {
		reasonsInvalid = append(reasonsInvalid, "claims are nil")
	} else {
		// device must be valid
		validationResponse, err := a.sf001DeviceValidator.Validate(&validator.ValidateRequest{
			Claims: request.Claims,
			Action: action.UpdateAllowedFields,
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

func (a *administrator) UpdateAllowedFields(request *administrator2.UpdateAllowedFieldsRequest) (*administrator2.UpdateAllowedFieldsResponse, error) {
	if err := a.ValidateUpdateAllowedFieldsRequest(request); err != nil {
		return nil, err
	}

	// retrieve the device
	deviceRetrieveResponse, err := a.sf001RecordHandler.Retrieve(&recordHandler.RetrieveRequest{
		Claims:     request.Claims,
		Identifier: id.Identifier{Id: request.SF001.Id},
	})
	if err != nil {
		return nil, exception.DeviceRetrieval{Reasons: []string{err.Error()}}
	}

	// update the allowed fields on the device

	// update the device
	_, err = a.sf001RecordHandler.Update(&recordHandler.UpdateRequest{
		Claims:     request.Claims,
		Identifier: id.Identifier{Id: request.SF001.Id},
		SF001:      deviceRetrieveResponse.SF001,
	})
	if err != nil {
		return nil, exception.DeviceUpdate{Reasons: []string{err.Error()}}
	}

	return &administrator2.UpdateAllowedFieldsResponse{
		SF001: deviceRetrieveResponse.SF001,
	}, nil
}
